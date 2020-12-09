package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Shopify/sarama"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/ads/internal/actions"
	"github.com/schoren/example-adserver/ads/internal/handlers"
	"github.com/schoren/example-adserver/ads/internal/platform/kafka"
	"github.com/schoren/example-adserver/ads/internal/platform/mysql"
	"github.com/schoren/example-adserver/pkg/config"
	"github.com/schoren/example-adserver/pkg/instrumentation"
	"github.com/schoren/example-adserver/pkg/retry"
)

type appConfig struct {
	AdserverBaseURL       string   `env:"ADSERVER_BASE_URL" validate:"nonzero"`
	DBDSN                 string   `env:"DB_DSN" validate:"nonzero"`
	SrvAddr               string   `env:"SRV_ADDR" validate:"nonzero"`
	KafkaBootstrapServers []string `env:"KAFKA_BOOTSTRAP_SERVERS" envSeparator:"," validate:"nonzero"`
}

func main() {
	cfg := appConfig{}
	config.MustReadFromEnv(&cfg)

	db := openDBConnection(cfg)
	defer db.Close()

	kafkaProducer := connectKafkaProducer(cfg)
	defer kafkaProducer.Close()

	router := mux.NewRouter().PathPrefix("/ads").Subrouter()
	setupHandlers(cfg, db, kafkaProducer, router)
	srv := setupHTTPServer(cfg, router)

	log.Printf("Starting server on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func setupHandlers(cfg appConfig, db *sql.DB, kafkaProducer sarama.SyncProducer, r *mux.Router) {
	kafkaNotifier := kafka.NewNotifier(kafkaProducer, config.KafkaTopicsAdUpdates)
	adsRepository := mysql.NewAdsRepository(db)

	inst := instrumentation.NewLogger(log.New(os.Stderr, "[Ads]", log.LstdFlags))

	create := handlers.NewCreate(actions.NewCreator(adsRepository, kafkaNotifier, inst), cfg.AdserverBaseURL)
	create.Register(r)

	update := handlers.NewUpdate(actions.NewUpdater(adsRepository, kafkaNotifier, inst))
	update.Register(r)

	listActive := handlers.NewListActive(actions.NewActiveLister(adsRepository, inst))
	listActive.Register(r)

}

func setupHTTPServer(cfg appConfig, r *mux.Router) *http.Server {

	return &http.Server{
		Handler: r,
		Addr:    cfg.SrvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
}

func openDBConnection(cfg appConfig) *sql.DB {
	db, err := sql.Open("mysql", cfg.DBDSN)
	if err != nil {
		panic(fmt.Errorf("cannot connect to mysql database: %w", err))
	}

	return db
}

func connectKafkaProducer(cfg appConfig) sarama.SyncProducer {
	config := sarama.NewConfig()

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true
	var kafkaProducer sarama.SyncProducer
	err := retry.Do(func() error {
		var err error
		kafkaProducer, err = sarama.NewSyncProducer(cfg.KafkaBootstrapServers, config)
		if err != nil {
			err = fmt.Errorf("cannot connect to kafka bootstrap servers: %v", err)
			log.Println(err)
			return err
		}
		log.Println("Connected to kafka cluster")
		return nil
	}, 5*time.Second, 5)

	if err != nil {
		panic(err)
	}

	return kafkaProducer
}
