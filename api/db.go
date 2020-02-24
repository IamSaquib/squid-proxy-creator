package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"

	// SQLite3 import
	_ "github.com/mattn/go-sqlite3"
)

// Config to strore conf
type Config struct {
	ID    string   `json:"id"`
	Peer  []string `json:"peers"`
	Host  string   `json:"host"`
	Port  string   `json:"port"`
	State int32    `json:"state"`
	Ts    string   `json:"ts"`
	TsMod string   `json:"ts_mod"`
}

// AddToDB to insert config into DB
func AddToDB(configuration Config) (*uuid.UUID, error) {
	db, err := sql.Open("sqlite3", "proxy.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	id := uuid.New()
	stmt, err := tx.Prepare("INSERT INTO proxy_config(id, peers, host, port, state) values (?, json(?), ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	mPeer, _ := json.Marshal(configuration.Peer)
	defer stmt.Close()
	_, err = stmt.Exec(
		id,
		string(mPeer),
		configuration.Host,
		configuration.Port,
		configuration.State,
	)
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return &id, nil
}

// ShowDB to list all config
func ShowDB() ([]Config, error) {
	var conf []Config
	db, err := sql.Open("sqlite3", "proxy.db")
	if err != nil {
		return conf, err
	}
	defer db.Close()
	rows, err := db.Query("select id, json_array(proxy_config.peers), host, port, state, ts, ts_mod from proxy_config")
	if err != nil {
		return conf, err
	}
	defer rows.Close()
	for rows.Next() {
		var confn Config
		var marshalledPeer string
		err = rows.Scan(&confn.ID, &marshalledPeer, &confn.Host, &confn.Port, &confn.State, &confn.Ts, &confn.TsMod)
		if err != nil {
			return conf, err
		}
		json.Unmarshal([]byte(marshalledPeer), &confn.Peer)
		conf = append(conf, Config{
			ID:    confn.ID,
			Peer:  confn.Peer,
			Host:  confn.Host,
			Port:  confn.Port,
			State: confn.State,
			Ts:    confn.Ts,
			TsMod: confn.TsMod,
		})
	}
	return conf, nil
}

// ShowByID to show details about a particular proxy
func ShowByID(id string) (Config, error) {
	var conf Config
	db, err := sql.Open("sqlite3", "proxy.db")
	if err != nil {
		return conf, err
	}
	defer db.Close()
	stmt, err := db.Prepare("select id, json_array(proxy_config.peers), host, port, state, ts, ts_mod from proxy_config where id = ?")
	if err != nil {
		return conf, err
	}
	defer stmt.Close()
	var nConf Config
	var marshalledPeer string
	err = stmt.QueryRow(id).Scan(&nConf.ID, &marshalledPeer, &nConf.Host, &nConf.Port, &nConf.State, &nConf.Ts, &nConf.TsMod)
	if err != nil {
		return conf, err
	}
	json.Unmarshal([]byte(marshalledPeer), &nConf.Peer)
	conf = Config{
		ID:    id,
		Peer:  nConf.Peer,
		Host:  nConf.Host,
		Port:  nConf.Port,
		State: nConf.State,
		Ts:    nConf.Ts,
		TsMod: nConf.TsMod,
	}

	return conf, nil
}

// UpdateDB to update a particular config by the ID
func UpdateDB(configuration Config) error {
	db, err := sql.Open("sqlite3", "proxy.db")
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE proxy_config set peers=?, ts_mod=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		configuration.Peer,
		time.Now(),
		configuration.ID,
	)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// DeleteDB to delete a particular row by ID
func DeleteDB(configuration Config) error {
	db, err := sql.Open("sqlite3", "proxy.db")
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("DELETE FROM proxy_config where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		configuration.ID,
	)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// GetPort method to get the next available port
func GetPort() (*string, error) {
	db, err := sql.Open("sqlite3", "proxy.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	rows, err := tx.Query("SELECT port_number from proxy_port where availability=1 LIMIT 1")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var pNum int
	for rows.Next() {
		err = rows.Scan(&pNum)
		if err != nil {
			return nil, err
		}
		stmt, err := tx.Prepare("UPDATE proxy_port SET availability=0 where port_number=?")
		if err != nil {
			return nil, err
		}
		defer stmt.Close()
		_, err = stmt.Exec(
			pNum,
		)
		if err != nil {
			return nil, err
		}
	}
	tx.Commit()
	fmt.Println("Port: %v", pNum)
	portNum := strconv.Itoa(int(pNum))
	return &portNum, nil
}
