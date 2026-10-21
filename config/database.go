package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ConnectDatabase() {
	dsn := fmt.Sprintf("host='%s' port='%s' user='%s' password='%s' dbname='%s' sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_USER", "agungrilo"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_NAME", "bioskop_db"),
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("gagal membuka koneksi database: ", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("gagal terhubung ke database: ", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS bioskop (
		id SERIAL PRIMARY KEY,
		nama VARCHAR(255) NOT NULL,
		lokasi VARCHAR(255) NOT NULL,
		rating REAL
	)`)
	if err != nil {
		log.Fatal("gagal membuat tabel bioskop: ", err)
	}

	DB = db
	fmt.Println("berhasil terhubung ke database")
}
