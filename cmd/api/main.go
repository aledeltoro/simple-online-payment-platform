package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aledeltoro/simple-online-payment-platform/cmd/api/handler"
	"github.com/aledeltoro/simple-online-payment-platform/internal/database/postgres"
	"github.com/aledeltoro/simple-online-payment-platform/internal/paymentprocessor/stripe"
	"github.com/aledeltoro/simple-online-payment-platform/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("load .env file failed: %s", err.Error())
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3000"
	}

	ctx := context.Background()

	database, err := postgres.Init(ctx)
	if err != nil {
		log.Fatalf("initialize database failed: %s \n", err.Error())
	}

	defer database.Close()

	paymentprocessor, err := stripe.New()
	if err != nil {
		log.Fatalf("initialize stripe payment processor failed: %s \n", err.Error())
	}

	onlinePaymentService := service.NewOnlinePaymentService(database, paymentprocessor)

	handler := handler.NewHandler(onlinePaymentService)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Route("/payments", func(r chi.Router) {
		r.Post("/", http.HandlerFunc(handler.HandleProcessPayment(ctx)))
		r.Get("/{id}", http.HandlerFunc(handler.HandleQueryPayment(ctx)))
		r.Post("/{id}/refunds", http.HandlerFunc(handler.HandleRefundPayment(ctx)))
	})

	fmt.Printf("Listening on port %s \n", port)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}
