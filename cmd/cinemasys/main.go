package main

import (
	"cinemasys/internal/database"
	"cinemasys/internal/router"
	"cinemasys/internal/server"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbConnectionString := os.Getenv("DATABASE_URL")
	if dbConnectionString == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	sqlxDB, err := sqlx.Connect("postgres", dbConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlxDB.Close()

	db := database.NewDatabase(sqlxDB)
	err = db.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := router.NewRouter(db)
	srv := server.NewServer("8888", r)
	srv.StartServer(*r)
}
