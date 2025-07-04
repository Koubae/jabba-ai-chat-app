package main

import (
	"github.com/Koubae/jabba-ai-chat-app/internal/chat-identity"
)

func main() {
	chat_identity.RunServer()
}

//import (
//	"database/sql"
//	"fmt"
//	"log"
//	"time"
//
//	_ "github.com/go-sql-driver/mysql"
//)

//
//func main() {
//	// Database connection string
//	// Format: username:password@tcp(host:port)/database_name
//	dsn := "root:admin@tcp(localhost:3306)/chat_identity"
//
//	// Open database connection
//	db, err := sql.Open("mysql", dsn)
//	if err != nil {
//		log.Fatal("Error opening database:", err)
//	}
//	defer db.Close()
//
//	// Maximum number of open connections to the database
//	db.SetMaxOpenConns(25)
//
//	// Maximum number of idle connections in the pool
//	db.SetMaxIdleConns(25)
//
//	// Maximum amount of time a connection may be reused
//	db.SetConnMaxLifetime(5 * time.Minute)
//
//	// Maximum amount of time a connection may be idle before being closed
//	db.SetConnMaxIdleTime(10 * time.Minute) // Optional, Go 1.15+
//
//	// Test the connection
//	err = db.Ping()
//	if err != nil {
//		log.Fatal("Error connecting to database:", err)
//	}
//
//	fmt.Println("Successfully connected to MySQL!")
//
//}
