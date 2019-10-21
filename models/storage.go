package models

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Storage - all methods for use DB
type Storage interface {
	CreateUser(*User)
	SetDataSearch(*User)
	GetAllUsers() []User
	SetIsSearching(*User)
	Close()
}

//MySQLStorage provider that can handle read/write from database
type MySQLStorage struct {
	con *sql.DB
}

// NewMysql will open db connection or return error
func NewMysql(host, user, password, dbname string) (*MySQLStorage, error) {
	if host == "" {
		log.Fatal("Empty host string, setup DB_HOST env")
		host = "localhost"
	}

	if user == "" {
		return nil, fmt.Errorf("Empty user string, setup DB_USER env")
	}

	if dbname == "" {
		return nil, fmt.Errorf("Empty dbname string, setup DB_DBNAME env")
	}

	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("Cannot open mysql connection: %v", err)
	}
	return &MySQLStorage{con: db}, nil
}

func (db *MySQLStorage) Close() {
	db.con.Close()
}
