package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kiRiLL3311/Currency/news/internal/config"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Get("DB_HOST"),
		config.Get("DB_PORT"),
		config.Get("DB_USER"),
		config.Get("DB_PASSWORD"),
		config.Get("DB_NAME"),
	)

	database, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.Ping(); err != nil {
		log.Fatal("DB not connected:", err)
	}

	log.Println("PostgreSQL connected (news)")
	return database
}
