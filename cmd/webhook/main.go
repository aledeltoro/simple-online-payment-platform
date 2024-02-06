package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aledeltoro/simple-online-payment-platform/cmd/webhook/handler"
	"github.com/aledeltoro/simple-online-payment-platform/internal/database/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("load .env file failed: %s", err.Error())
	}

	port := os.Getenv("WEBHOOKS_PORT")
	if port == "" {
		port = "3001"
	}

	ctx := context.Background()

	database, err := postgres.Init(ctx)
	if err != nil {
		log.Fatalf("initialize database failed: %s \n", err.Error())
	}

	defer database.Close()

	handler := handler.NewHandler(database)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Post("/payments/{provider}/events", http.HandlerFunc(handler.HandlePaymentEvents(ctx)))

	fmt.Printf("Listening on port %s \n", port)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
