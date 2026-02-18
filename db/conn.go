// ...existing code...
package db

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func ConnectDB() (*sql.DB, error) {
	host := getenv("DB_HOST", "localhost")
	port := getenvInt("DB_PORT", 5433)
	user := getenv("DB_USER", "admin")
	password := getenv("DB_PASSWORD", "admin")
	dbname := getenv("DB_NAME", "postgres")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("Connected to " + dbname)
	return db, nil
}
