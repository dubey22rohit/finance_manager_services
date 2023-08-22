package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type BrokerConfig struct {
	Rabbit *amqp.Connection
}

func main() {
	var app BrokerConfig

	flag.Parse()

	conn, err := connRabbit()
	if err != nil {
		log.Println("error connecting to RabbitMQ", err)
		os.Exit(1)
	}
	defer conn.Close()

	app.Rabbit = conn

	log.Printf("starting broker service on port %s\n", webPort)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connRabbit() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not ready yet....")
			counts++
		} else {
			log.Println("connected to RabbitMQ!")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off....")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
