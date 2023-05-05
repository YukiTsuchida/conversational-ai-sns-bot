package main

import (
	"log"
	"net/http"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/controller/http/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

func main() {
	// Configを読み込む

	db := newEntClient()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Post("/accounts/twitter", handler.RegisterTwitterAccountHandler(db))
	http.ListenAndServe(":8080", r)
}

func newEntClient() *ent.Client {
	client, err := ent.Open("postgres", "postgresql://admin:admin@postgresql:5432/db?sslmode=disable")
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}
	return client
}
