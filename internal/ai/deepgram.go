package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

func (s Service) TranscribeAudio(ctx context.Context, audio []byte, contentType string) (string, error) {
	if s.deepgramAPIKey == "" {
		return "", errors.New("missing DEEPGRAM_API_KEY")
	}

	request, err := s.newDeepgramRequest(ctx, audio, contentType)
	if err != nil {
		return "", err
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	return parseDeepgramResponse(response)
}

func (s Service) newDeepgramRequest(ctx context.Context, audio []byte, contentType string) (*http.Request, error) {
	query := url.Values{}
	query.Set("model", s.deepgramModel)
	query.Set("smart_format", "true")
	query.Set("detect_language", "true")
	endpoint := "https://api.deepgram.com/v1/listen?" + query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(audio))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Token "+s.deepgramAPIKey)
	request.Header.Set("Content-Type", fallbackContentType(contentType))
	return request, nil
}

func parseDeepgramResponse(response *http.Response) (string, error) {
	if response.StatusCode >= http.StatusBadRequest {
		return "", errors.New("deepgram transcription failed")
	}

	var payload deepgramResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}

	return firstTranscript(payload)
}

func fallbackContentType(contentType string) string {
	if contentType == "" {
		return "audio/webm"
	}

	return contentType
}

func firstTranscript(payload deepgramResponse) (string, error) {
	if len(payload.Results.Channels) == 0 {
		return "", errors.New("deepgram returned no channels")
	}

	alternatives := payload.Results.Channels[0].Alternatives
	if len(alternatives) == 0 || alternatives[0].Transcript == "" {
		return "", errors.New("deepgram returned an empty transcript")
	}

	return alternatives[0].Transcript, nil
}
