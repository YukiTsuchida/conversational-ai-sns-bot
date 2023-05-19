package main

import (
	"log"
	"net/http"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/config"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/ent"
	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/lib/pq"
)

func main() {
	db := newEntClient()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Post("/conversations/{id}/send", handlers.SendConversationHandler(db))
	http.ListenAndServe(":8080", r)
}

func newEntClient() *ent.Client {
	dsn := "postgresql://" + config.POSTGRES_USER() + ":" + config.POSTGRES_PASSWORD() + "@" + config.POSTGRES_HOST() + ":" + config.POSTGRES_PORT() + "/" + config.POSTGRES_DB() + "?sslmode=disable"
	client, err := ent.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("creating client: %v", err)
	}
	return client
}
