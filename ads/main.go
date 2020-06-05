package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/ads/internal/handlers"
	"github.com/schoren/example-adserver/ads/internal/platform/kafka"
	"github.com/schoren/example-adserver/ads/internal/platform/mysql"
)

func main() {

	dbDSN := os.Getenv("DB_DSN")
	if dbDSN == "" {
		panic(fmt.Errorf("DB_DSN not provided"))
	}

	srvAddr := os.Getenv("SRV_ADDR")
	if srvAddr == "" {
		panic(fmt.Errorf("SRV_ADDR not provided"))
	}

	kafkaBootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstrapServers == "" {
		panic(fmt.Errorf("SRV_ADDR not provided"))
	}

	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		panic(fmt.Errorf("cannot connect to mysql database: %w", err))
	}
	defer db.Close()

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	// Ugly fix for docker-compose start order: just wait a few secs for kafka to be ready
	time.Sleep(10 * time.Second)
	log.Println("Waited enough, try to connect to", strings.Split(kafkaBootstrapServers, ","))
	kafkaProducer, err := sarama.NewSyncProducer(strings.Split(kafkaBootstrapServers, ","), config)
	if err != nil {
		panic(fmt.Errorf("cannot connect to kafka bootstrap servers: %w", err))
	}
	defer kafkaProducer.Close()

	adsRepository := mysql.NewAdsRepository(db)
	kafkaNotifier := kafka.NewNotifier(kafkaProducer)

	handlers.CreateCommand = &commands.Create{
		Persister: adsRepository,
		Notifier:  kafkaNotifier,
	}

	handlers.UpdateCommand = &commands.Update{
		Persister: adsRepository,
		Notifier:  kafkaNotifier,
	}

	router := mux.NewRouter()
	handlers.ConfigureRouter(router.PathPrefix("/ads").Subrouter())

	srv := &http.Server{
		Handler: router,
		Addr:    srvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}

	log.Printf("Starting server on %s", srvAddr)
	log.Fatal(srv.ListenAndServe())
}
