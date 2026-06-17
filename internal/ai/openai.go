package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func (s Service) GenerateResearch(ctx context.Context, transcript string) (ResearchResult, error) {
	if s.openAIAPIKey == "" {
		return ResearchResult{}, errors.New("missing OPENAI_API_KEY")
	}

	request, err := s.newOpenAIRequest(ctx, transcript)
	if err != nil {
		return ResearchResult{}, err
	}

	response, err := s.httpClient.Do(request)
	if err != nil {
		return ResearchResult{}, err
	}
	defer response.Body.Close()

	return parseOpenAIResponse(response)
}

func (s Service) newOpenAIRequest(ctx context.Context, transcript string) (*http.Request, error) {
	body, err := json.Marshal(buildOpenAIPayload(s.openAIModel, transcript))
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+s.openAIAPIKey)
	request.Header.Set("Content-Type", "application/json")
	return request, nil
}

func parseOpenAIResponse(response *http.Response) (ResearchResult, error) {
	if response.StatusCode >= http.StatusBadRequest {
		return ResearchResult{}, errors.New("openai generation failed")
	}

	var payload openAIResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return ResearchResult{}, err
	}

	text := firstOutputText(payload)
	if text == "" {
		return ResearchResult{}, errors.New("openai returned empty output")
	}

	var result ResearchResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return ResearchResult{}, err
	}

	return result, nil
}

func firstOutputText(payload openAIResponse) string {
	if payload.OutputText != "" {
		return payload.OutputText
	}

	for _, item := range payload.Output {
		for _, content := range item.Content {
			if content.Text != "" {
				return content.Text
			}
		}
	}

	return ""
}

func buildOpenAIPayload(model string, transcript string) map[string]any {
	return map[string]any{
		"model": model,
		"tools": []map[string]any{{"type": "web_search_preview"}},
		"input": []map[string]any{
			buildMessage("system", systemPrompt()),
			buildMessage("user", transcript),
		},
		"text": map[string]any{
			"format": map[string]any{
				"type":   "json_schema",
				"name":   "jarvis_research_response",
				"strict": true,
				"schema": researchSchema(),
			},
		},
	}
}

func buildMessage(role string, text string) map[string]any {
	return map[string]any{
		"role": role,
		"content": []map[string]string{{
			"type": "input_text",
			"text": text,
		}},
	}
}

func systemPrompt() string {
	return "You are Jarvis, an AI Chief of Staff. Research the request using web search when useful. Return concise, opinionated business analysis. Do not use filler."
}

func researchSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title":   map[string]string{"type": "string"},
			"summary": map[string]string{"type": "string"},
			"key_findings": map[string]any{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
			"recommendation": map[string]string{"type": "string"},
			"image_queries": map[string]any{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
		},
		"required":             []string{"title", "summary", "key_findings", "recommendation", "image_queries"},
		"additionalProperties": false,
	}
}
