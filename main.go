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
	fmt.Println("Insert some data to users table")
	if _, err := db.Exec(
		context.Background(),
		"INSERT INTO users (id, name) VALUES ($1, $2)",
		5, "John",
	); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	fmt.Println("Select all rows from users table")
	if rows, err := db.Query(
		context.Background(),
		"SELECT id, name FROM users ORDER BY id ASC",
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

	fmt.Println("Update inside transaction")
	if err := db.Transaction(context.Background(), func(ctx context.Context, tx *gosql.Tx) error {
		if _, err := tx.Exec(ctx, "UPDATE users SET name=$1 WHERE id=$2", "John", 1); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, "UPDATE users SET name=$1 WHERE id=$2", "Alice", 5); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	fmt.Println("Select all rows from users again")
	if rows, err := db.Query(
		context.Background(),
		"SELECT id, name FROM users ORDER BY id ASC",
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
