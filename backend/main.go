package main

import (
	"fmt"
	"log"

	"github.com/bukharney/giga-chat/configs"
	"github.com/bukharney/giga-chat/database"
	"github.com/bukharney/giga-chat/server"
)

func main() {
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
	fmt.Println("Server is running on port 8080")
}
