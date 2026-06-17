package api

type ResearchRequest struct {
	Transcript string `json:"transcript"`
}

type ResearchResponse struct {
	Action         string   `json:"action"`
	Transcript     string   `json:"transcript"`
	Title          string   `json:"title"`
	Summary        string   `json:"summary"`
	KeyFindings    []string `json:"keyFindings"`
	Recommendation string   `json:"recommendation"`
	Images         []string `json:"images"`
}
