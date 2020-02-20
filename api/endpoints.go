// API to perform CRUD for proxy servers

package api

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

// AddToConf function to append into /etc/squid/squid.conf
func AddToConf(path string) error {
	fs, err := os.OpenFile("squid.conf",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fs.Close()
	if _, err := fs.WriteString("\ninclude /squid/acl-" + path + ".conf"); err != nil {
		return err
	}
	return nil
}

// createACL function to generate ACL command for the given list IPs
func createACL(id string) error {
	fs, err := os.Create("squid/acl-" + id + ".conf")
	if err != nil {
		return err
	}
	defer fs.Close()
	var aclCmd bytes.Buffer
	aclCmd.WriteString("acl " + id + " dstdomain \"squid/iplist-" + id + ".acl\"\nhttp_access allow " + id)
	if _, err = fs.Write(aclCmd.Bytes()); err != nil {
		log.Fatal(err)
	}
	return nil
}

// CreateProxy to create user defined proxy
func CreateProxy(w http.ResponseWriter, r *http.Request) {
	var config Config
	_ = json.NewDecoder(r.Body).Decode(&config)
	id, err := AddToDB(config)
	if err != nil {
		log.Fatal(err)
	}
	path := "squid/iplist-" + id.String() + ".acl"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var ipList bytes.Buffer
	for _, ip := range config.Peer.Ips {
		ipList.WriteString(ip + "\n")
	}
	fmt.Println(ipList.String())
	if _, err = f.Write(ipList.Bytes()); err != nil {
		log.Fatal(err)
	}
	err = createACL(id.String())
	if err != nil {
		log.Fatal(err)
	}
	err = AddToConf(id.String())
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
	params := mux.Vars(r)
	id := params["id"]
	config, err := ShowByID(id)
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
	err := os.Remove("squid/" + strings.Replace(config.ID, " ", "", -1) + ".conf")
	if err != nil {
		log.Fatal(err)
	}
	path := "squid/" + strings.Replace(config.ID, " ", "", -1) + ".conf"
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	content := &bytes.Buffer{}
	gob.NewEncoder(content).Encode(config.Peer)
	contentSlice := content.Bytes()
	if _, err = f.Write(contentSlice); err != nil {
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
	// err := os.Remove("squid/" + strings.Replace(config.ID, " ", "", -1) + ".conf")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode("Deleted")
}
