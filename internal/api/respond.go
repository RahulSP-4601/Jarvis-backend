package api

import (
	"context"

	"jarvis-backend/internal/ai"
)

func (s *Server) respondToTranscript(ctx context.Context, transcript string) (ResearchResponse, error) {
	return s.respondToResearchRequest(ctx, ResearchRequest{Transcript: transcript})
}

func (s *Server) respondToResearchRequest(ctx context.Context, request ResearchRequest) (ResearchResponse, error) {
	if isCloseCommand(request.Transcript) {
		return ResearchResponse{
			Action:       "hide_overlay",
			Transcript:   request.Transcript,
			Title:        "Jarvis hidden",
			Summary:      "Closing the overlay.",
			SpokenAnswer: "Alright. I'll step back.",
		}, nil
	}

	result, err := s.ai.GenerateResearch(ctx, ai.ResearchRequest{
		Transcript:      request.Transcript,
		Locale:          request.Locale,
		AcceptedLocales: request.AcceptedLocales,
	})
	if err != nil {
		return ResearchResponse{}, err
	}

	return buildResponse(request.Transcript, result), nil
}

func buildResponse(transcript string, result ai.ResearchResult) ResearchResponse {
	return ResearchResponse{
		Action:          "respond",
		Transcript:      transcript,
		Title:           result.Title,
		Summary:         result.Summary,
		KeyFindings:     result.KeyFindings,
		Recommendation:  result.Recommendation,
		Images:          toImageURLs(result.ImageQueries),
		SpokenAnswer:    result.SpokenAnswer,
		FollowUpPrompts: result.FollowUpPrompts,
	}
}
