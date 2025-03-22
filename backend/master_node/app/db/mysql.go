package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func DBConfig() *sql.DB {
	mySqlDB := os.Getenv("MYSQL_DATABASE")
	mySqlUser := os.Getenv("MYSQL_USER")
	mySqlPassword := os.Getenv("MYSQL_PASSWORD")
	mySqlHost := os.Getenv("MYSQL_HOST")

	dsn := mySqlUser + ":" + mySqlPassword + "@tcp(" + mySqlHost + ":3306)/" + mySqlDB + "?parseTime=true"
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database: ", err.Error())
		return nil
	}

	fmt.Println("Successfully connected to the database!")
	return db
}
