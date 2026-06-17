package ai

type ResearchResult struct {
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	KeyFindings    []string `json:"key_findings"`
	Recommendation string   `json:"recommendation"`
	ImageQueries   []string `json:"image_queries"`
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
