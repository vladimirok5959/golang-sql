package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/vladimirok5959/golang-sql/gosql"
)

func main() {
	// Get temp file name
	f, err := ioutil.TempFile("", "go-sqlite-")
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
	f.Close()

	// Set migration directory
	migrationsDir, err := filepath.Abs("./db/migrations")
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	// Open DB connection, SQLite is used as example
	// You can use here MySQL or PostgreSQL, just change dbURL
	db, err := gosql.Open("sqlite://"+f.Name(), migrationsDir, true)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	// DB struct here ./db/migrations/20220527233113_test_migration.sql
	// Insert some data to users table:
	if _, err := db.Exec(
		context.Background(),
		"INSERT INTO users (id, name) VALUES ($1, $2)",
		5, "john",
	); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	// Select all rows from users table:
	if rows, err := db.Query(
		context.Background(),
		"SELECT id, name FROM users ORDER BY id DESC",
	); err == nil {
		type rowStruct struct {
			ID   int64
			Name string
		}
		defer rows.Close()
		for rows.Next() {
			var row rowStruct
			if err := rows.Scan(&row.ID, &row.Name); err != nil {
				panic(fmt.Sprintf("%s", err))
			}
			fmt.Printf("ID: %d, Name: %s\n", row.ID, row.Name)
		}
		if err := rows.Err(); err != nil {
			panic(fmt.Sprintf("%s", err))
		}
	} else {
		panic(fmt.Sprintf("%s", err))
	}

	// Close DB connection
	if err := db.Close(); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
