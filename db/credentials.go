package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbhost     = "localhost"
	dbuser     = "prime"
	dbpassword = "belindat2014"
	dbname     = "avy"
	dbport     = "5432"
)

func DBconnect() *sql.DB {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbuser, dbpassword, dbhost, dbport, dbname)

	dbconn, err := sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	db = dbconn

	return db
}
