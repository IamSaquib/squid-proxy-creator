// API to perform CRUD for proxy servers

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// AppendToFile which appends content to a file
func AppendToFile(content string) error {
	f, err := os.OpenFile("squid.conf",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	log.Println("Command :: " + content)
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		return err
	}
	return nil
}

// CreateProxy to create user defined proxy
func CreateProxy(w http.ResponseWriter, r *http.Request) {
	var config Config
	_ = json.NewDecoder(r.Body).Decode(&config)
	if err := AddToDB(config); err != nil {
		log.Fatal(err)
	}
	f, err := os.Create("squid/squid1.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	content := []byte(config.Config)
	if _, err = f.Write(content); err != nil {
		log.Fatal(err)
	}
}

// ShowProxy to show all proxies available for user
func ShowProxy(w http.ResponseWriter, r *http.Request) {
	var conf []Config
	conf, err := ShowDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(conf)
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(conf)
}

// ShowProxyByID to show details about proxy by ID
func ShowProxyByID(w http.ResponseWriter, r *http.Request) {
	var conf Config
	_ = json.NewDecoder(r.Body).Decode(&conf)
	fmt.Print(conf)
	config, err := ShowByID(conf.ID)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(config)
}

// UpdateProxy to update proxy based on config id
func UpdateProxy(w http.ResponseWriter, r *http.Request) {
	var config Config
	_ = json.NewDecoder(r.Body).Decode(&config)
	if err := UpdateDB(config); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode("Updated")
}
