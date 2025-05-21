package main

import (
	"github.com/TMateusz1/go-3rd-devs/internal/ai"
	"github.com/TMateusz1/go-3rd-devs/internal/ai/option"
	"github.com/TMateusz1/go-3rd-devs/internal/middleware"
	"github.com/TMateusz1/go-3rd-devs/internal/websearch"
	"github.com/TMateusz1/go-3rd-devs/websearch/handler"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"net/http"
)

func main() {
	as, err := ai.NewOpenaiService(option.WithBaseModel(ai.OpenRouterModelGPT4oMini))
	if err != nil {
		log.Fatalln(err)
	}

	ws, err := websearch.NewService(as)
	if err != nil {
		log.Fatalln(err)
	}
	wh := handler.NewWebSearchHandler(as, ws)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/websearch", middleware.LogMiddleware(wh.Handle))

	s := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on: ", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatalln(err)
	}
}
