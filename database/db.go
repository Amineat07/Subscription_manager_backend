package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitiateDataBase() *pgxpool.Pool{
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env")
	}

	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	fmt.Println("Connecting to DB:", os.Getenv("DB_DBNAME"))

	pool, err := pgxpool.New(context.Background(), psqlInfo)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	fmt.Println("Database successfully connected!")
	return pool
}