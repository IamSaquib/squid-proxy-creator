package main

import (
	"context"
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/squid-proxy-creator/api"
)

func main() {

	// if err := api.ShowDB(); err != nil {
	// 	log.Fatal(err)
	// }
	// var conf &api.Config
	// conf.userID = ""
	// conf := api.Config{
	// 	UserID:    "bab775c9-96a5-459b-a823-29a3faca7d39",
	// 	Config:    "acl src 192.12.343.42",
	// 	ProxyName: "My first proxy",
	// }
	// if err := api.CreateProxy(conf); err != nil {
	// 	log.Fatal("Ran into ", err)
	// }

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	// Add your routes as needed
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	r.HandleFunc("/create-proxy", api.CreateProxy)
	r.HandleFunc("/show-proxy", api.ShowProxy)
	r.HandleFunc("/show-proxy-id", api.ShowProxyByID)
	r.HandleFunc("/update-proxy", api.UpdateProxy)
	log.Println("Running server on :1406")
	srv := &http.Server{
		Addr: "0.0.0.0:1406",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}
