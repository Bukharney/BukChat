package main

import (
	"log"

	"github.com/bukharney/giga-chat/configs"
	"github.com/bukharney/giga-chat/database"
	"github.com/bukharney/giga-chat/server"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	cfg := configs.NewConfigs()

	db, err := database.NewPostgreSQL(cfg)
	if err != nil {
		log.Fatal(err)
	}

	srv := server.NewServer(db, cfg)
	err = srv.Run()
	if err != nil {
		log.Fatal(err)
	}
}
