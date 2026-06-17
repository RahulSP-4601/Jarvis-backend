package api

import (
	"context"

	"jarvis-backend/internal/ai"
)

func (s *Server) respondToTranscript(ctx context.Context, transcript string) (ResearchResponse, error) {
	if isCloseCommand(transcript) {
		return ResearchResponse{
			Action:     "hide_overlay",
			Transcript: transcript,
			Title:      "Jarvis hidden",
			Summary:    "Closing the overlay.",
		}, nil
	}

	result, err := s.ai.GenerateResearch(ctx, transcript)
	if err != nil {
		return ResearchResponse{}, err
	}

	return buildResponse(transcript, result), nil
}

func buildResponse(transcript string, result ai.ResearchResult) ResearchResponse {
	return ResearchResponse{
		Action:         "respond",
		Transcript:     transcript,
		Title:          result.Title,
		Summary:        result.Summary,
		KeyFindings:    result.KeyFindings,
		Recommendation: result.Recommendation,
		Images:         toImageURLs(result.ImageQueries),
	}
}
