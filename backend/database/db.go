package database

import (
	"errors"
	"log"

	"github.com/bukharney/giga-chat/configs"
	"github.com/bukharney/giga-chat/utils"
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
	"name" VARCHAR(255) UNIQUE NOT NULL,
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "messages" (
	"id" SERIAL PRIMARY KEY,
	"room_id" INTEGER NOT NULL,
	"user_id" INTEGER NOT NULL,
	"message" TEXT NOT NULL,
	"created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS "users_rooms" (
	"user_id" INTEGER NOT NULL,
	"room_id" INTEGER NOT NULL,
	PRIMARY KEY ("user_id", "room_id")
);

CREATE TABLE IF NOT EXISTS "friends" (
	"from_user_id" INTEGER NOT NULL,
	"to_user_id" INTEGER NOT NULL,
	"status" INTEGER NOT NULL,
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
		return nil, errors.New(err.Error())
	}

	db.MustExec(schema)

	log.Println("Connected to PostgreSQL")
	return db, nil
}
