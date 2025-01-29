package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vladimirok5959/golang-sql/gosql"
)

func main() {
	// Get temp file name
	f, err := os.CreateTemp("", "go-sqlite-")
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
	db, err := gosql.Open("sqlite://"+f.Name(), migrationsDir, false, true)
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	db.SetConnMaxLifetime(time.Minute * 60)
	db.SetMaxIdleConns(8)
	db.SetMaxOpenConns(8)

	// DB struct here ./db/migrations/20220527233113_test_migration.sql
	fmt.Println("Inserting some data to users table")
	if _, err := db.Exec(
		context.Background(),
		"INSERT INTO users (id, name) VALUES ($1, $2)",
		5, "John",
	); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	fmt.Println("Selecting all rows from users table")
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

	fmt.Println("Updating inside transaction")
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

	fmt.Println("Selecting all rows from users again")
	if err := db.Each(
		context.Background(),
		"SELECT id, name FROM users ORDER BY id ASC",
		func(ctx context.Context, rows *gosql.Rows) error {
			var row struct {
				ID   int64
				Name string
			}
			if err := rows.Scans(&row); err != nil {
				return err
			}
			fmt.Printf("ID: %d, Name: %s\n", row.ID, row.Name)
			return nil
		},
	); err != nil {
		panic(fmt.Sprintf("%s", err))
	}

	fmt.Println("Selecting specific user with ID: 5")
	var row struct {
		ID   int64
		Name string
	}
	err = db.QueryRow(context.Background(), "SELECT id, name FROM users WHERE id=$1", 5).Scans(&row)
	if err != nil && err != sql.ErrNoRows {
		panic(fmt.Sprintf("%s", err))
	} else {
		if err != sql.ErrNoRows {
			fmt.Printf("ID: %d, Name: %s\n", row.ID, row.Name)
		} else {
			fmt.Printf("Record not found\n")
		}
	}

	// Close DB connection
	if err := db.Close(); err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}
