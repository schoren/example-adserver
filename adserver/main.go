package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
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
)

var topics = []string{"ad-updates"}

func main() {
	srvAddr := os.Getenv("SRV_ADDR")
	if srvAddr == "" {
		panic(fmt.Errorf("SRV_ADDR not provided"))
	}

	kafkaBootstrapServers := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if kafkaBootstrapServers == "" {
		panic(fmt.Errorf("SRV_ADDR not provided"))
	}

	adServiceBaseURL := os.Getenv("AD_SERVICE_BASE_URL")
	if adServiceBaseURL == "" {
		panic(fmt.Errorf("AD_SERVICE_BASE_URL not provided"))
	}

	// Ugly fix for docker-compose start order: just wait a few secs for kafka to be ready
	time.Sleep(15 * time.Second)
	log.Println("Waited enough, try to connect to", strings.Split(kafkaBootstrapServers, ","))

	adLister := rest.NewAdLister(adServiceBaseURL)
	adStore := memory.NewAdStore()

	err := adstore.Warmup(adStore, adLister)
	if err != nil {
		panic(fmt.Errorf("Cannot Warmup AdStore: %w", err))
	}

	serveCommand := &commands.NewServe(adStore)
	updateAdCommand := &commands.UpdateAdCommand{AdStore: adStore}

	handlers.ServeCommand = serveCommand

	router := mux.NewRouter()
	handlers.ConfigureRouter(router.PathPrefix("/").Subrouter())

	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	adUpdater := kafka.NewAdUpdater(updateAdCommand)

	ctx := context.Background()
	client, err := sarama.NewConsumerGroup(strings.Split(kafkaBootstrapServers, ","), "adserver-"+uuid.New().String(), config)
	if err != nil {
		panic(fmt.Errorf("Error creating Kafka consumer group client: %w", err))
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := client.Consume(ctx, topics, adUpdater); err != nil {
				panic(fmt.Errorf("Error from kafka consumer: %w", err))
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()

	<-adUpdater.Ready // Await till the consumer has been set up
	log.Println("Kafka consumer up and running!...")

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
