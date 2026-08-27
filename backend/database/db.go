package database

import (
	"fmt"
	"log"

	"github.com/bukharney/bukchat/configs"
	"github.com/bukharney/bukchat/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS "users" (
	"id" SERIAL PRIMARY KEY,
	"username" VARCHAR(255) UNIQUE NOT NULL,
	"email" VARCHAR(255) UNIQUE NOT NULL,
	"password" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "rooms" (
	"id" SERIAL PRIMARY KEY,
	"name" VARCHAR(255) NOT NULL,
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "messages" (
	"id" SERIAL PRIMARY KEY,
	"room_id" INTEGER NOT NULL REFERENCES rooms(id),
	"user_id" INTEGER NOT NULL REFERENCES users(id),
	"message" TEXT NOT NULL,
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "users_rooms" (
	"user_id" INTEGER NOT NULL REFERENCES users(id),
	"room_id" INTEGER NOT NULL REFERENCES rooms(id),
	PRIMARY KEY ("user_id", "room_id")
);

CREATE TABLE IF NOT EXISTS "friends" (
	"from_user_id" INTEGER NOT NULL REFERENCES users(id),
	"to_user_id" INTEGER NOT NULL REFERENCES users(id),
	"status" INTEGER NOT NULL,
	"room_id" INTEGER REFERENCES rooms(id),
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY ("from_user_id", "to_user_id")
);
`

func NewPostgreSQL(cfg *configs.Configs) (*sqlx.DB, error) {
	connectionUrl, err := utils.ConnectionUrlBuilder("postgres", cfg)
	if err != nil {
		return nil, err
	}

	log.Println("Connecting to PostgreSQL")

	log.Println(connectionUrl)
	db, err := sqlx.Connect("postgres", connectionUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgresql: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to execute database schema: %w", err)
	}

	log.Println("Connected to PostgreSQL")
	return db, nil
}

