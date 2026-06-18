package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"runtime"
	"strings"

	"jarvis-backend/internal/config"
)

func (s *Server) handleHealth(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleResearch(writer http.ResponseWriter, request *http.Request) {
	var payload ResearchRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	response, err := s.respondToResearchRequest(request.Context(), payload)
	if err != nil {
		writeJSON(writer, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func (s *Server) handleBootstrap(writer http.ResponseWriter, request *http.Request) {
	var payload BootstrapRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	accessToken, err := bearerToken(request)
	if err != nil {
		writeJSON(writer, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	user, err := s.auth.GetUser(request.Context(), accessToken)
	if err != nil {
		writeJSON(writer, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(writer, http.StatusOK, BootstrapResponse{
		User: UserProfile{
			ID:    user.ID,
			Email: user.Email,
			Name:  resolveUserName(user.UserMetadata.FullName, user.UserMetadata.Name, user.Email),
		},
		Runtime: RuntimeConfig{
			APIBaseURL:          buildRuntimeAPIBaseURL(s.config),
			NativeWakeEnabled:   s.config.PicovoiceAccessKey != "",
			NativeWakeProvider:  "picovoice",
			NativeWakeAccessKey: s.config.PicovoiceAccessKey,
		},
		Features: FeatureFlags{
			ResearchEnabled:    true,
			GoogleOAuthEnabled: true,
		},
	})
}

func (s *Server) handleVoiceCommand(writer http.ResponseWriter, request *http.Request) {
	audio, contentType, err := readAudioUpload(request)
	if err != nil {
		writeJSON(writer, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	transcript, err := s.ai.TranscribeAudio(request.Context(), audio, contentType)
	if err != nil {
		writeJSON(writer, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	response, err := s.respondToResearchRequest(request.Context(), ResearchRequest{Transcript: transcript})
	if err != nil {
		writeJSON(writer, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(writer, http.StatusOK, response)
}

func readAudioUpload(request *http.Request) ([]byte, string, error) {
	file, header, err := request.FormFile("audio")
	if err != nil {
		return nil, "", errors.New("missing audio upload")
	}
	defer file.Close()

	payload, err := io.ReadAll(file)
	if err != nil {
		return nil, "", errors.New("failed to read audio upload")
	}

	return payload, header.Header.Get("Content-Type"), nil
}

func writeJSON(writer http.ResponseWriter, statusCode int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(payload)
}

func bearerToken(request *http.Request) (string, error) {
	authorization := strings.TrimSpace(request.Header.Get("Authorization"))
	if authorization == "" {
		return "", errors.New("missing authorization header")
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return "", errors.New("invalid authorization header")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
	if token == "" {
		return "", errors.New("missing bearer token")
	}

	return token, nil
}

func resolveUserName(fullName string, name string, email string) string {
	if fullName != "" {
		return fullName
	}

	if name != "" {
		return name
	}

	return strings.Split(email, "@")[0]
}

func buildRuntimeAPIBaseURL(cfg config.Config) string {
	if cfg.PublicAPIBaseURL != "" {
		return cfg.PublicAPIBaseURL
	}

	if cfg.Port == "" {
		return ""
	}

	host := "http://localhost:"
	if runtime.GOOS == "windows" {
		return host + cfg.Port
	}

	return host + cfg.Port
}
