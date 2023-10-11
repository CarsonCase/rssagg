package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/CarsonCase/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load()

	fmt.Printf("Hello World\n")

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Port not found. Please fill out ENV file")
	}

	sqlString := os.Getenv("DB_URL")

	if sqlString == "" {
		log.Fatal("sql string not found. Please fill out ENV file")
	}

	connection, err := sql.Open("postgres", sqlString)

	if err != nil {
		log.Fatal("Error connecting to db")
	}

	apiCfg := apiConfig{
		DB: database.New(connection),
	}

	go startScraping(database.New(connection), 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	routerV1 := chi.NewRouter()

	routerV1.Get("/ready", handlerReadiness)
	routerV1.Get("/error", handlerErr)
	routerV1.Post("/users", apiCfg.handlerCreateUser)
	routerV1.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	routerV1.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	routerV1.Get("/feeds", apiCfg.handlerGetFeeds)
	routerV1.Post("/feedFollows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))

	router.Mount("/v1", routerV1)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v\n", portString)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(portString)
}
