package main

import (
	"database/sql"
	"log"

	"github.com/kitessafikadu/kevent/internal/database"
	"github.com/kitessafikadu/kevent/internal/env"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

type application struct {
	port     int
	jwtSecret string
	models database.Models
}

func main(){
	db,err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port: env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models: models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}