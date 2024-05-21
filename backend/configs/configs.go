package configs

import (
	"log"
	"os"
)

type Configs struct {
	PostgreSQL PostgreSQL
	App        Gin
}

type Gin struct {
	Host string
	Port string
}

type PostgreSQL struct {
	Host     string
	Port     string
	Protocol string
	Username string
	Password string
	Database string
	SSLMode  string
}

func NewConfigs() *Configs {
	return &Configs{
		PostgreSQL: PostgreSQL{
			Host:     MustGetenv("POSTGRES_HOST"),
			Port:     MustGetenv("POSTGRES_PORT"),
			Protocol: "tcp",
			Username: MustGetenv("POSTGRES_USER"),
			Password: MustGetenv("POSTGRES_PASSWORD"),
			Database: MustGetenv("POSTGRES_DB"),
			SSLMode:  "disable",
		},
		App: Gin{
			Host: "localhost",
			Port: "8080",
		},
	}
}

func MustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing env var %s", key)
	}
	return v
}
