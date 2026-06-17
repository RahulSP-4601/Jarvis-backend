package ai

import (
	"net/http"
	"time"

	"jarvis-backend/internal/config"
)

type Service struct {
	openAIAPIKey   string
	openAIModel    string
	deepgramAPIKey string
	deepgramModel  string
	httpClient     *http.Client
}

func NewService(cfg config.Config) Service {
	return Service{
		openAIAPIKey:   cfg.OpenAIAPIKey,
		openAIModel:    cfg.OpenAIModel,
		deepgramAPIKey: cfg.DeepgramAPIKey,
		deepgramModel:  cfg.DeepgramModel,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}
