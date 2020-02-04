package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type IP struct {
	Address string `json:"address"`
}

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

func createCommand(command string) string {
	return "\nacl myclient src " + command + "\nhttp_access allow myclient"
}

// Method to add users IP Address
func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var ip IP
	_ = json.NewDecoder(r.Body).Decode(&ip)
	address := ip.Address
	command := createCommand(address)
	if err := AppendToFile(command); err != nil {
		log.Fatalln(err)
	}
	json.NewEncoder(w).Encode(&ip)
}
