package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	DB       *sql.DB
	Models   any
	port     int
	env      string
	dbconfig struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  int
	}
}

func main() {
	log.Println("starting auth service")

	var cfg Config

	flag.IntVar(&cfg.port, "port", 4000, "Auth service server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")

	flag.StringVar(&cfg.dbconfig.dsn, "auth-sb-dsn", "", "Auth PostgresSQL DB dns")

	flag.IntVar(&cfg.dbconfig.maxIdleConns, "db-max-idle-conns", 25, "DB max idle connections")
	flag.IntVar(&cfg.dbconfig.maxOpenConns, "db-max-open-conns", 25, "DB max open connections")
	flag.IntVar(&cfg.dbconfig.maxIdleTime, "db-max-idle-time", 15, "DB max idle time")

	db, err := openDB(cfg)
	if err != nil {
		log.Fatalf("unable to start auth db %v", err)
	}

	app := Config{
		DB: db,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic("error starting server: ", err)
	}
}

func openDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dbconfig.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.dbconfig.maxIdleConns)
	db.SetMaxOpenConns(cfg.dbconfig.maxOpenConns)
	db.SetConnMaxIdleTime(time.Duration(cfg.dbconfig.maxIdleTime))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
