// API to perform CRUD for proxy servers

package api

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/gorilla/mux"
)

var mutex sync.Mutex

// Whitelist to store whitelist IPs
type Whitelist struct {
	IP string `json:"ip"`
}

// AddToConf function to append into /etc/squid/squid.conf
func AddToConf(path string, wg *sync.WaitGroup) error {
	mutex.Lock()
	fs, err := os.OpenFile("squid.conf",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fs.Close()
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(currentDir)

	if _, err := fs.WriteString("\ninclude " + currentDir + "/squid/acl-" + path + ".conf"); err != nil {
		return err
	}
	mutex.Unlock()
	wg.Done()
	return nil
}

// createACL function to generate ACL command for the given list IPs
func createACL(id string, port string) error {
	fs, err := os.Create("squid/acl-" + id + ".conf")
	if err != nil {
		return err
	}
	defer fs.Close()
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(currentDir)
	var aclBody bytes.Buffer

	aclBody.WriteString("http_port " + port + " name=port-" + id)
	aclBody.WriteString("\nacl " + id + " dstdomain \"/home/iamsaquib/Dev/my_open_source/squid_proxy_balancer/squid/iplist-" + id + ".acl\"")
	aclBody.WriteString("\nacl user-" + id + " myportname port-" + id)
	aclBody.WriteString("\nhttp_access allow user-" + id + " " + id)
	if _, err = fs.Write(aclBody.Bytes()); err != nil {
		log.Fatal(err)
	}
	return nil
}

// selectNewPort function to return available port to add for a particular config
// func

// CreateProxy to create user defined proxy
func CreateProxy(w http.ResponseWriter, r *http.Request) {
	var config Config
	_ = json.NewDecoder(r.Body).Decode(&config)
	port, err := GetPort()
	if err != nil {
		log.Fatal(err)
	}
	config.Port = *port
	id, err := AddToDB(config)
	if err != nil {
		log.Fatal(err)
	}
	path := "squid/iplist-" + id.String() + ".acl"

	fmt.Println("Writing to file with mutex lock")
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var ipList bytes.Buffer
	for _, ip := range config.Peer {
		ipList.WriteString(ip + "\n")
	}
	fmt.Println(ipList.String())
	if _, err = f.Write(ipList.Bytes()); err != nil {
		log.Fatal(err)
	}

	err = createACL(id.String(), config.Port)
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go AddToConf(id.String(), &wg)
	wg.Wait()
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode("Created")
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
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(currentDir)
	err = os.Remove("/home/iamsaquib/Dev/my_open_source/squid_proxy_balancer/squid/acl-" + strings.Replace(config.ID, " ", "", -1) + ".conf")
	if err != nil {
		log.Fatal(err)
	}
	err = os.Remove("/home/iamsaquib/Dev/my_open_source/squid_proxy_balancer/squid/iplist-" + strings.Replace(config.ID, " ", "", -1) + ".acl")
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode("Deleted")
}

// AddIPWhitelist to add white list IPs
func AddIPWhitelist(w http.ResponseWriter, r *http.Request) {
	var wl Whitelist
	_ = json.NewDecoder(r.Body).Decode(&wl)
	fs, err := os.OpenFile("whitelist.txt", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()
	if _, err = fs.WriteString(wl.IP + "\n"); err != nil {
		log.Fatal(err)
	}
}

// ShowWhitelist to show the list of IPs added
func ShowWhitelist(w http.ResponseWriter, r *http.Request) {
	fs, err := os.OpenFile("whitelist.txt", os.O_RDONLY, 0400)
	if err != nil {
		log.Fatal(err)
	}
	defer fs.Close()
	// scanner := bufio.NewScanner(fs)
	ips, err := ioutil.ReadAll(fs)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "Application/json")
	json.NewEncoder(w).Encode(string(ips))
}

// RemoveWhitelist to remove a particular whitelist IP
func RemoveWhitelist(w http.ResponseWriter, r *http.Request) {
	var wl Whitelist
	_ = json.NewDecoder(r.Body).Decode(&wl)
	fs, err := ioutil.ReadFile("whitelist.txt")
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("(?m)[\r\n]+^.*" + wl.IP + ".*$")
	res := re.ReplaceAllString(string(fs), "")
	err = ioutil.WriteFile("whitelist.txt", []byte(res), 0644)
}
