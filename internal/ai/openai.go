package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func (s Service) GenerateResearch(ctx context.Context, request ResearchRequest) (ResearchResult, error) {
	if s.openAIAPIKey == "" {
		return ResearchResult{}, errors.New("missing OPENAI_API_KEY")
	}

	httpRequest, err := s.newOpenAIRequest(ctx, request)
	if err != nil {
		return ResearchResult{}, err
	}

	response, err := s.httpClient.Do(httpRequest)
	if err != nil {
		return ResearchResult{}, err
	}
	defer response.Body.Close()

	return parseOpenAIResponse(response)
}

func (s Service) newOpenAIRequest(ctx context.Context, request ResearchRequest) (*http.Request, error) {
	body, err := json.Marshal(buildOpenAIPayload(s.openAIModel, request))
	if err != nil {
		return nil, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.openai.com/v1/responses", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpRequest.Header.Set("Authorization", "Bearer "+s.openAIAPIKey)
	httpRequest.Header.Set("Content-Type", "application/json")
	return httpRequest, nil
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

func buildOpenAIPayload(model string, request ResearchRequest) map[string]any {
	return map[string]any{
		"model": model,
		"tools": []map[string]any{{"type": "web_search_preview"}},
		"input": []map[string]any{
			buildMessage("system", systemPrompt(request)),
			buildMessage("user", request.Transcript),
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

func systemPrompt(request ResearchRequest) string {
	languageInstruction := buildLanguageInstruction(request)

	return "You are Jarvis, an AI Chief of Staff focused on research. Research the request using web search when useful. Speak like a mature, sharp business partner. Be concise, opinionated, and practical. Challenge weak assumptions when needed. Prioritize useful judgment over generic explanation. Return a spoken answer first, then a structured brief with findings, recommendation, follow-up prompts, and image queries when they add value. " + languageInstruction
}

func buildLanguageInstruction(request ResearchRequest) string {
	if request.Locale == "" && len(request.AcceptedLocales) == 0 {
		return "Reply in the user's language when it is clear from the request."
	}

	acceptedLocales := strings.Join(request.AcceptedLocales, ", ")
	if acceptedLocales == "" {
		return "Prefer replying in locale " + request.Locale + " when practical."
	}

	return "Reply in the same language the user speaks whenever practical. Preferred locale: " + request.Locale + ". Accepted locales: " + acceptedLocales + "."
}

func researchSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"title":         map[string]string{"type": "string"},
			"summary":       map[string]string{"type": "string"},
			"spoken_answer": map[string]string{"type": "string"},
			"key_findings": map[string]any{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
			"recommendation": map[string]string{"type": "string"},
			"follow_up_prompts": map[string]any{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
			"image_queries": map[string]any{
				"type":  "array",
				"items": map[string]string{"type": "string"},
			},
		},
		"required":             []string{"title", "summary", "spoken_answer", "key_findings", "recommendation", "follow_up_prompts", "image_queries"},
		"additionalProperties": false,
	}
}
