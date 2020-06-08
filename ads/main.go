package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/caarlos0/env/v6"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/ads/internal/handlers"
	"github.com/schoren/example-adserver/ads/internal/platform/kafka"
	"github.com/schoren/example-adserver/ads/internal/platform/mysql"
	"gopkg.in/validator.v2"
)

type config struct {
	AdserverBaseURL       string   `env:"ADSERVER_BASE_URL" validate:"nonzero"`
	DBDSN                 string   `env:"DB_DSN" validate:"nonzero"`
	SrvAddr               string   `env:"SRV_ADDR" validate:"nonzero"`
	KafkaBootstrapServers []string `env:"KAFKA_BOOTSTRAP_SERVERS" envSeparator:"," validate:"nonzero"`
}

func readEnvConfig(cfg interface{}) {
	if err := env.Parse(cfg); err != nil {
		panic(fmt.Errorf("cannot parse env config: %w", err))
	}

	if errs := validator.Validate(cfg); errs != nil {
		panic(fmt.Errorf("invalid env config: %w", errs))
	}
}

func main() {
	// Ugly fix for docker-compose start order: just wait a few secs for kafka to be ready
	time.Sleep(10 * time.Second)
	log.Println("Waited enough, try to connect to", cfg.KafkaBootstrapServers)

	cfg := config{}

	readEnvConfig(&cfg)

	db := openDBConnection(cfg)
	defer db.Close()

	kafkaProducer := connectKafkaProducer(cfg)
	defer kafkaProducer.Close()

	setupHandlers(cfg, db, kafkaProducer)
	srv := setupHTTPServer(cfg)

	log.Printf("Starting server on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

	kafkaNotifier := kafka.NewNotifier(kafkaProducer, config.KafkaTopicsAdUpdates)
	adsRepository := mysql.NewAdsRepository(db)

	handlers.AdServerBaseURL = cfg.AdserverBaseURL
	handlers.CreateCommand = commands.NewCreate(adsRepository, kafkaNotifier)
	handlers.UpdateCommand = commands.NewUpdate(adsRepository, kafkaNotifier)
	handlers.ListActiveCommand = commands.NewListActive(adsRepository)

}

func setupHTTPServer(cfg config) *http.Server {
	router := mux.NewRouter()
	handlers.ConfigureRouter(router.PathPrefix("/ads").Subrouter())

	srv := &http.Server{
		Handler: router,
		Addr:    cfg.SrvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
}

func openDBConnection(cfg config) *sql.DB {
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		panic(fmt.Errorf("cannot connect to mysql database: %w", err))
	}

	return db
}

func connectKafkaProducer(cfg config) sarama.SyncProducer {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	kafkaProducer, err := sarama.NewSyncProducer(cfg.KafkaBootstrapServers, config)

	if err != nil {
		panic(fmt.Errorf("cannot connect to kafka bootstrap servers: %w", err))
	}

	return kafkaProducer
}
