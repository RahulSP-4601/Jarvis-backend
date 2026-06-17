package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

	response, err := s.respondToTranscript(request.Context(), payload.Transcript)
	if err != nil {
		writeJSON(writer, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(writer, http.StatusOK, response)
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

	response, err := s.respondToTranscript(request.Context(), transcript)
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
