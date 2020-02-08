package api

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// Config to structure conf
type Config struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Config    string `json:"config"`
	ProxyName string `json:"proxy_name"`
}

// AddToDB to insert config into DB
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

// ShowDB to list all config
func ShowDB() ([]Config, error) {
	var conf []Config
	db, err := sql.Open("sqlite3", "proxy_config.db")
	if err != nil {
		return conf, err
	}
	defer db.Close()
	rows, err := db.Query("select * from proxy_config")
	if err != nil {
		return conf, err
	}
	defer rows.Close()
	for rows.Next() {
		var userID string
		var config string
		var proxyName string
		var id string
		err = rows.Scan(&id, &userID, &config, &proxyName)
		if err != nil {
			return conf, err
		}
		conf = append(conf, Config{
			ID:        id,
			UserID:    userID,
			Config:    config,
			ProxyName: proxyName,
		})
	}
	return conf, nil
}

// ShowByID to show details about a particular proxy
func ShowByID(id string) (Config, error) {
	var conf Config
	db, err := sql.Open("sqlite3", "proxy_config.db")
	if err != nil {
		return conf, err
	}
	defer db.Close()
	stmt, err := db.Prepare("select * from proxy_config where id = ?")
	if err != nil {
		return conf, err
	}
	defer stmt.Close()
	var userID string
	var config string
	var proxyName string
	var ID string
	err = stmt.QueryRow(id).Scan(&ID, &userID, &config, &proxyName)
	if err != nil {
		return conf, err
	}
	conf = Config{
		ID:        id,
		UserID:    userID,
		Config:    config,
		ProxyName: proxyName,
	}

	return conf, nil
}

// UpdateDB to update a particular config by the ID
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
	stmt, err := tx.Prepare("UPDATE proxy_config set config=? where id=?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		configuration.Config,
		configuration.ID,
	)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}
