package config

import "os"

type Config struct {
	Port               string
	AllowedOrigin      string
	OpenAIAPIKey       string
	OpenAIModel        string
	DeepgramAPIKey     string
	DeepgramModel      string
	SupabaseURL        string
	SupabaseAnonKey    string
	PicovoiceAccessKey string
	PublicAPIBaseURL   string
}

func Load() Config {
	return Config{
		Port:               getenv("PORT", "8080"),
		AllowedOrigin:      getenv("ALLOWED_ORIGIN", "http://localhost:5173"),
		OpenAIAPIKey:       os.Getenv("OPENAI_API_KEY"),
		OpenAIModel:        getenv("OPENAI_MODEL", "gpt-4.1-mini"),
		DeepgramAPIKey:     os.Getenv("DEEPGRAM_API_KEY"),
		DeepgramModel:      getenv("DEEPGRAM_MODEL", "nova-3"),
		SupabaseURL:        os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:    os.Getenv("SUPABASE_ANON_KEY"),
		PicovoiceAccessKey: os.Getenv("PICOVOICE_ACCESS_KEY"),
		PublicAPIBaseURL:   getenv("PUBLIC_API_BASE_URL", ""),
	}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
