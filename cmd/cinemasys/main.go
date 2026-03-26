package main

import (
	"cinemasys/internal/database"
	"cinemasys/internal/router"
	"cinemasys/internal/server"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
)

func main() {
	godotenv.Load()
	// redisClient := database.NewRedisClient(os.Getenv("REDIS_URL"))
	db, err := sqlx.Connect("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	database := database.NewDatabase(db)
	if err := database.InitDB(); err != nil {
		panic(err)
	}

	router := router.NewRouter(database)

	server := server.NewServer(os.Getenv("PORT"), router)
	server.Initialize()
}
