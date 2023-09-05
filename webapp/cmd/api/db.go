package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (app *Application) ConnectToDB() (*sql.DB, error) {
	db, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to Postgres")
	return db, nil
}
