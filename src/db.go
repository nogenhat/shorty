package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/lib/pq"
)

// DB tools ====================================================================

var db *sql.DB

func confirmDbReset() {
	msg := fmt.Sprintf("Do you really want to reset DB '%s'?", config.Database.Database)
	if !confirm(msg) {
		halt()
	}
}

func initDBSchema() {
	schema := `
		DROP TABLE IF EXISTS urls;

		CREATE SEQUENCE urls_id_seq;
		CREATE SEQUENCE urls_code_seq START 0 MINVALUE 0;

		CREATE TABLE IF NOT EXISTS urls (
			id int4 PRIMARY KEY NOT NULL DEFAULT nextval('urls_id_seq'),
			url varchar(255) NOT NULL,
			code varchar(32) NOT NULL,
			open_count int4 NOT NULL DEFAULT 0
		);

		CREATE UNIQUE INDEX urls_code_ind ON urls (code);

		ALTER SEQUENCE urls_id_seq OWNED BY urls.id;
		ALTER SEQUENCE urls_code_seq OWNED BY urls.code;`

	if _, err := db.Exec(schema); err != nil {
		checkErr(err, "Can't init DB schema")
	}
}

func initDB(resetConfirmation bool) {
	var err error

	connectionString := fmt.Sprintf(
		"host='%s' dbname='%s' user='%s' password='%s' sslmode=disable",
		config.Database.Host,
		config.Database.Database,
		config.Database.User,
		config.Database.Password,
	)

	db, err = sql.Open("postgres", connectionString)
	checkErr(err, "Can't connect to DB")

	db.SetMaxOpenConns(config.Database.MaxOpenConnections)
	db.SetMaxIdleConns(config.Database.MaxIdleConnections)

	if config.Database.InitSchema {
		if resetConfirmation {
			confirmDbReset()
		}
		initDBSchema()
	}
}

func closeDB() {
	db.Close()
}

// end of DB tools

// DB queries ==================================================================

const codeBase = 13330

func createUrl(url string) (code string, err error) {
	var codeSeq int64

	err = db.QueryRow("SELECT nextval('urls_code_seq')").Scan(&codeSeq)
	if err != nil {
		return
	}

	code = strconv.FormatInt(codeBase+codeSeq, 36)

	_, err = db.Exec("INSERT INTO urls (url, code) VALUES ($1, $2)", url, code)
	if err != nil {
		return
	}

	return
}

func getUrl(code string) (url string, err error) {
	err = db.QueryRow(
		"SELECT url FROM urls WHERE code = $1 LIMIT 1",
		code,
	).Scan(&url)

	return
}

func hitRedirect(code string) (err error) {
	_, err = db.Exec(
		"UPDATE urls SET open_count = open_count + 1 WHERE code = $1",
		code,
	)

	return
}

func getOpenCount(code string) (count int64, err error) {
	err = db.QueryRow(
		"SELECT open_count FROM urls WHERE code = $1 LIMIT 1",
		code,
	).Scan(&count)

	return
}

// end of DB queries
