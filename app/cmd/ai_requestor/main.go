package main

import (
	"net/http"

	"github.com/YukiTsuchida/conversational-ai-sns-bot/app/http/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Post("/openai_chat_gpt", handlers.ProxyOpenAIChatGPTHandler())
	http.ListenAndServe(":8080", r)
}
