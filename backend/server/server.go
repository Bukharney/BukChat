package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bukharney/bukchat/configs"
	_repo "github.com/bukharney/bukchat/modules/repositories"
	"github.com/bukharney/bukchat/server/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	App *gin.Engine
	Cfg *configs.Configs
	DB  *sqlx.DB
}

func NewServer(db *sqlx.DB, cfg *configs.Configs) *Server {
	return &Server{
		App: gin.Default(),
		DB:  db,
		Cfg: cfg,
	}
}

func (s *Server) Run() error {
	s.App.Use(cors.New(
		cors.Config{
			AllowOrigins:     []string{"*"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			AllowCredentials: true,
		},
	))

	hub := ws.NewHub()
	go hub.Run()

	err := s.MapHandlers(hub)
	if err != nil {
		return fmt.Errorf("failed to map handlers: %w", err)
	}

	chatRepo := _repo.NewChatRepo(s.DB)
	s.App.GET("/ws/:roomId", func(c *gin.Context) {
		ws.ServeWS(c, hub, chatRepo)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: s.App,
	}

	go func() {
		slog.Info("Starting HTTP server", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP Server Listener Error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutdown signal received, shutting down HTTP server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("HTTP Server Forced Shutdown Error", "error", err)
		return err
	}

	slog.Info("Server exited cleanly")
	return nil
}
