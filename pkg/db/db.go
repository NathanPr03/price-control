package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func ConnectToDb() (*sql.DB, error) {
	_ = godotenv.Load("../../.env.development.local")

	var (
		host     = os.Getenv("POSTGRES_HOST")
		port     = 5432
		user     = os.Getenv("POSTGRES_USER")
		password = os.Getenv("POSTGRES_PASSWORD")
		dbname   = os.Getenv("POSTGRES_DATABASE")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return &sql.DB{}, err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return &sql.DB{}, err
	}
	fmt.Println("Successfully connected!")

	return db, nil
}
