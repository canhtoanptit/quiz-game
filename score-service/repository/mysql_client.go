package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type MysqlClient struct {
	db *sql.DB
}

func NewMysqlClient(mysqlConf string) *MysqlClient {
	db, err := sql.Open("mysql", mysqlConf)
	if err != nil {
		panic(err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	} else {
		fmt.Println("Successfully connected to MySQL!")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &MysqlClient{
		db: db,
	}
}
