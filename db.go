package main

import (
	_ "github.com/lib/pq"
)

// func InitDB() *sql.DB {
// 	connStr := "postgres://user:password@localhost:5432/ecofoot?sslmode=disable"
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Failed to connect database:", err)
// 	}

// 	if err := db.Ping(); err != nil {
// 		log.Fatal("Database ping error:", err)
// 	}

// 	return db
// }
