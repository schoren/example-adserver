package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/schoren/example-adserver/ads/internal/commands"
	"github.com/schoren/example-adserver/ads/internal/handlers"
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

	db, err := sql.Open("mysql", dbDSN)
	if err != nil {
		panic(fmt.Errorf("cannot connect to mysql database: %w", err))
	}
	defer db.Close()

	router := mux.NewRouter()

	adsRepository := mysql.NewAdsRepository(db)

	handlers.CreateCommand = &commands.Create{
		Persister: adsRepository,
	}

	handlers.UpdateCommand = &commands.Update{
		Persister: adsRepository,
	}

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
