package api

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Config to structure conf
type Config struct {
	Id        string
	UserID    string
	Config    string
	ProxyName string
}

func AddToDB(configuration Config) error {
	db, err := sql.Open("sqlite3", "proxy_config.db")
	if err != nil {
		return err
	}

	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO proxy_config(user_id, config, proxy_name) values (?,?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		configuration.UserID,
		configuration.Config,
		configuration.ProxyName,
	)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func ShowDB() error {
	db, err := sql.Open("sqlite3", "proxy_config.db")
	if err != nil {
		return err
	}
	defer db.Close()
	rows, err := db.Query("select id, user_id, config, proxy_name from proxy_config")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var userID string
		var config string
		var proxyName string
		var id string
		err = rows.Scan(&id, &userID, &config, &proxyName)
		if err != nil {
			return err
		}
		fmt.Println(id, userID, config, proxyName)
	}
	return nil
}

func UpdateDB(configuration Config) error {
	db, err := sql.Open("sqlite3", "proxy_config.db")
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("UPDATE proxy_config set config=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		configuration.Config,
	)
	if err != nil {
		return err
	}
	return nil
}
