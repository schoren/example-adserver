package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/adserver/internal/adstore"
	"github.com/schoren/example-adserver/adserver/internal/commands"
	"github.com/schoren/example-adserver/adserver/internal/handlers"
	"github.com/schoren/example-adserver/adserver/internal/platform/kafka"
	"github.com/schoren/example-adserver/adserver/internal/platform/memory"
	"github.com/schoren/example-adserver/adserver/internal/platform/rest"
	"github.com/schoren/example-adserver/pkg/config"
	"github.com/schoren/example-adserver/pkg/retry"
)

type appConfig struct {
	AdServiceBaseURL      string   `env:"AD_SERVICE_BASE_URL" validate:"nonzero"`
	SrvAddr               string   `env:"SRV_ADDR" validate:"nonzero"`
	KafkaBootstrapServers []string `env:"KAFKA_BOOTSTRAP_SERVERS" envSeparator:"," validate:"nonzero"`
}

func main() {
	cfg := appConfig{}
	config.MustReadFromEnv(&cfg)

	adStore := createAdStore(cfg)

	setupHandlers(cfg, adStore)

	adUpdater := setupAdUpdater(cfg, adStore)

	<-adUpdater.Ready // Await till the consumer has been set up
	log.Println("Kafka consumer up and running!...")

	srv := setupHTTPServer(cfg)

	log.Printf("Starting server on %s", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}

func createAdStore(cfg appConfig) adstore.GetSetter {
	adLister := rest.NewAdLister(cfg.AdServiceBaseURL)
	adStore := memory.NewAdStore()

	err := retry.Do(func() error {
		err := adstore.Warmup(adStore, adLister)
		if err != nil {
			err := fmt.Errorf("Cannot Warmup AdStore: %v", err)
			log.Println(err)
			return err
		}
		return nil
	}, 5*time.Second, 5)
	if err != nil {
		panic(err)
	}

	return adStore
}

func setupAdUpdater(cfg appConfig, adStore adstore.GetSetter) *kafka.AdUpdater {
	adUpdater := kafka.NewAdUpdater(commands.NewUpdateAd(adStore))

	consumerGroup := "adserver-" + uuid.New().String()

	kfkConfig := sarama.NewConfig()
	kfkConfig.Version = sarama.MaxVersion

	var client sarama.ConsumerGroup
	ctx := context.Background()

	err := retry.Do(func() error {
		var err error
		client, err = sarama.NewConsumerGroup(cfg.KafkaBootstrapServers, consumerGroup, kfkConfig)
		if err != nil {
			err = fmt.Errorf("Error creating Kafka consumer group client: %v", err)
			log.Println(err)
			return err
		}
		log.Println("Connected to kafka cluster")
		return nil
	}, 5*time.Second, 5)

	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, []string{config.KafkaTopicsAdUpdates}, adUpdater); err != nil {
				panic(fmt.Errorf("Error from kafka consumer: %w", err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return adUpdater
}

func setupHandlers(cfg appConfig, adStore adstore.GetSetter) {
	handlers.ServeCommand = commands.NewServe(adStore)
}

func setupHTTPServer(cfg appConfig) *http.Server {
	router := mux.NewRouter()
	handlers.ConfigureRouter(router.PathPrefix("/").Subrouter())

	return &http.Server{
		Handler: router,
		Addr:    cfg.SrvAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 2 * time.Second,
		ReadTimeout:  2 * time.Second,
	}
}
