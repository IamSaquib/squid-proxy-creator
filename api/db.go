package api

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	// SQLite3 import
	_ "github.com/mattn/go-sqlite3"
)

// Config to structure conf
type Config struct {
	ID     string `json:"id"`
	Peer   Peer   `json:"peers"`
	Server string `json:"server"`
	State  int32  `json:"state"`
	Ts     string `json:"ts"`
	TsMod  string `json:"ts_mod"`
}

// Peer Struct to structure Ips
type Peer struct {
	Ips []string `json:"ips"`
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
	stmt, err := tx.Prepare("INSERT INTO proxy_config(id, peers, server, state) values (?, json(?), ?, ?)")
	if err != nil {
		return nil, err
	}
	mPeer, _ := json.Marshal(configuration.Peer)
	defer stmt.Close()
	_, err = stmt.Exec(
		id,
		string(mPeer),
		configuration.Server,
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
	rows, err := db.Query("select id, json_extract(proxy_config.peers, '$'), server, state, ts, ts_mod from proxy_config")
	if err != nil {
		return conf, err
	}
	defer rows.Close()
	for rows.Next() {
		var confn Config
		var marshalledPeer string
		err = rows.Scan(&confn.ID, &marshalledPeer, &confn.Server, &confn.State, &confn.Ts, &confn.TsMod)
		if err != nil {
			return conf, err
		}
		json.Unmarshal([]byte(marshalledPeer), &confn.Peer)
		conf = append(conf, Config{
			ID:     confn.ID,
			Peer:   confn.Peer,
			Server: confn.Server,
			State:  confn.State,
			Ts:     confn.Ts,
			TsMod:  confn.TsMod,
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
	stmt, err := db.Prepare("select id, json_extract(proxy_config.peers, '$'), server, state, ts, ts_mod from proxy_config where id = ?")
	if err != nil {
		return conf, err
	}
	defer stmt.Close()
	var nConf Config
	var marshalledPeer string
	err = stmt.QueryRow(id).Scan(&nConf.ID, &marshalledPeer, &nConf.Server, &nConf.State, &nConf.Ts, &nConf.TsMod)
	if err != nil {
		return conf, err
	}
	json.Unmarshal([]byte(marshalledPeer), &nConf.Peer)
	conf = Config{
		ID:     id,
		Peer:   nConf.Peer,
		Server: nConf.Server,
		State:  nConf.State,
		Ts:     nConf.Ts,
		TsMod:  nConf.TsMod,
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
