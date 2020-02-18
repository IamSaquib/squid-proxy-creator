package api

import (
	"database/sql"
	"time"

	// SQLite3 import
	_ "github.com/mattn/go-sqlite3"
)

// Config to structure conf
type Config struct {
	ID     string   `json:"id"`
	Peer   []string `json:"peers"`
	Server string   `json:"server"`
	State  int32    `json:"state"`
	Ts     string   `json:"ts"`
	TsMod  string   `json:"ts_mod"`
}

// AddToDB to insert config into DB
func AddToDB(configuration Config) error {
	db, err := sql.Open("sqlite3", "proxy.db")
	if err != nil {
		return err
	}

	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO proxy_config(peers, server, state) values (json('?'), ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		configuration.Peer,
		configuration.Server,
		configuration.State,
	)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

// ShowDB to list all config
func ShowDB() ([]Config, error) {
	var conf []Config
	db, err := sql.Open("sqlite3", "proxy.db")
	if err != nil {
		return conf, err
	}
	defer db.Close()
	rows, err := db.Query("select id, json_extract(proxy_config.peers, '$.ip'), server, state, ts, ts_mod from proxy_config")
	if err != nil {
		return conf, err
	}
	defer rows.Close()
	for rows.Next() {
		var confn Config
		err = rows.Scan(&confn.ID, &confn.Peer, &confn.Server, &confn.State, &confn.Ts, &confn.TsMod)
		if err != nil {
			return conf, err
		}
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
	stmt, err := db.Prepare("select * from proxy_config where id = ?")
	if err != nil {
		return conf, err
	}
	defer stmt.Close()
	var server string
	var peer []string
	var ID string
	var state int32
	err = stmt.QueryRow(id).Scan(&ID, &peer, &server, &state)
	if err != nil {
		return conf, err
	}
	conf = Config{
		ID:     id,
		Peer:   peer,
		Server: server,
		State:  state,
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
