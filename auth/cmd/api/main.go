package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dubey22rohit/finance_manager_services/auth/cmd/data"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

type AppConfig struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("starting auth service")

	db := connectToDB()
	if db == nil {
		log.Fatalf("unable to start auth db")
		return
	}
	defer db.Close()

	app := AppConfig{
		DB:     db,
		Models: data.New(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic("error starting server: ", err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	db.SetConnMaxIdleTime(time.Duration(15))

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	var counts int32

	dsn := os.Getenv("DSN")
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres is not ready yet....")
			counts++
		} else {
			log.Println("connected to postgres!")
			return conn
		}

		if counts > 10 {
			log.Println("error connecting to postgres", err)
			return nil
		}

		log.Println("Backing off for 4 seconds....")
		time.Sleep(4 * time.Second)
		continue
	}
}
