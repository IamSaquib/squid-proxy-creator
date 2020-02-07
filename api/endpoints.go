// API to perform CRUD for proxy servers

package api

import (
	"encoding/json"
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

// func ReplaceConf()

// Method to add users IP Address
//func AddUser(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	var ip IP
//	_ = json.NewDecoder(r.Body).Decode(&ip)
//	address := ip.Address
//	command := createCommand(address)
//	if err := AppendToFile(command); err != nil {
//		log.Fatalln(err)
//	}
//	json.NewEncoder(w).Encode(&ip)
//}
