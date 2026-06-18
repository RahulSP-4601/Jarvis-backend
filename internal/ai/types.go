package ai

type ResearchRequest struct {
	Transcript      string
	Locale          string
	AcceptedLocales []string
}

type ResearchResult struct {
	Title           string   `json:"title"`
	Summary         string   `json:"summary"`
	SpokenAnswer    string   `json:"spoken_answer"`
	KeyFindings     []string `json:"key_findings"`
	Recommendation  string   `json:"recommendation"`
	FollowUpPrompts []string `json:"follow_up_prompts"`
	ImageQueries    []string `json:"image_queries"`
}

type deepgramResponse struct {
	Results struct {
		Channels []struct {
			Alternatives []struct {
				Transcript string `json:"transcript"`
			} `json:"alternatives"`
		} `json:"channels"`
	} `json:"results"`
}

type openAIResponse struct {
	OutputText string `json:"output_text"`
	Output     []struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
}
