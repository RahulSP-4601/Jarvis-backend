package main

import (
	"log"
	"net/http"

	"jarvis-backend/internal/api"
	"jarvis-backend/internal/config"
)

func main() {
	cfg := config.Load()
	server := api.NewServer(cfg)

	log.Printf("jarvis-backend listening on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, server.Handler()); err != nil {
		log.Fatal(err)
	}
}
