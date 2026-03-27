package main

import (
	"cinemasys/internal/cache"
	"cinemasys/internal/database"
	"cinemasys/internal/router"
	"cinemasys/internal/server"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	dbURL := getEnv("DATABASE_URL", "postgres://cinemasys:secretpassword@localhost:5432/cinemasys?sslmode=disable")
	port := getEnv("PORT", "8080")

	db, err := sqlx.Connect("pgx", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to PostgreSQL")

	database := database.NewDatabase(db)
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database initialized successfully")

	redisAddr := getEnv("REDIS_ADDR", "localhost:6379")
	redisClient := cache.NewCache(redisAddr)
	if err := redisClient.HealthCheck(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()
	router := router.NewRouter(database, redisClient)
	srv := server.NewServer(port, router)

	log.Printf("Starting server on port %s", port)
	if err := srv.Initialize(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
