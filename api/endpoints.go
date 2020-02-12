// API to perform CRUD for proxy servers

package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// AddToConf function to append into /etc/squid/squid.conf
func AddToConf(path string) error {
	fs, err := os.OpenFile("squid.conf",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fs.Close()
	if _, err := fs.WriteString("\ninclude /" + path); err != nil {
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
	path := "squid/" + strings.Replace(config.ProxyName, " ", "", -1) + ".conf"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	content := []byte(config.Config)
	if _, err = f.Write(content); err != nil {
		log.Fatal(err)
	}
	err = AddToConf(path)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode("Created ")
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
	err := os.Remove("squid/" + strings.Replace(config.ProxyName, " ", "", -1) + ".conf")
	if err != nil {
		log.Fatal(err)
	}
	path := "squid/" + strings.Replace(config.ProxyName, " ", "", -1) + ".conf"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	content := []byte(config.Config)
	if _, err = f.Write(content); err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode("Updated")
}

// DeleteProxy to delete proxy based on config id
func DeleteProxy(w http.ResponseWriter, r *http.Request) {
	var config Config
	_ = json.NewDecoder(r.Body).Decode(&config)
	if err := DeleteDB(config); err != nil {
		log.Fatal(err)
	}
	err := os.Remove("squid/" + strings.Replace(config.ProxyName, " ", "", -1) + ".conf")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode("Deleted")
}
