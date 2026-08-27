package main

import (
	"log/slog"
	"os"

	"github.com/bukharney/bukchat/configs"
	"github.com/bukharney/bukchat/database"
	"github.com/bukharney/bukchat/pkg/logger"
	"github.com/bukharney/bukchat/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	logger.InitLogger()

	cfg := configs.NewConfigs()

	db, err := database.NewPostgreSQL(cfg)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	srv := server.NewServer(db, cfg)
	err = srv.Run()
	if err != nil {
		slog.Error("Server runtime error", "error", err)
		os.Exit(1)
	}
}
