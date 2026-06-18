package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"jarvis-backend/internal/config"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	UserMetadata struct {
		FullName string `json:"full_name"`
		Name     string `json:"name"`
	} `json:"user_metadata"`
}

func NewClient(cfg config.Config) Client {
	return Client{
		baseURL:    strings.TrimRight(cfg.SupabaseURL, "/"),
		apiKey:     cfg.SupabaseAnonKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c Client) GetUser(ctx context.Context, accessToken string) (User, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return User{}, errors.New("supabase auth is not configured")
	}

	if accessToken == "" {
		return User{}, errors.New("missing access token")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/auth/v1/user", nil)
	if err != nil {
		return User{}, err
	}

	request.Header.Set("Authorization", "Bearer "+accessToken)
	request.Header.Set("apikey", c.apiKey)

	response, err := c.httpClient.Do(request)
	if err != nil {
		return User{}, err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		return User{}, errors.New("invalid auth token")
	}

	var user User
	if err := json.NewDecoder(response.Body).Decode(&user); err != nil {
		return User{}, err
	}

	if user.ID == "" {
		return User{}, errors.New("authenticated user missing id")
	}

	return user, nil
}
