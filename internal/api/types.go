package api

type ResearchRequest struct {
	Transcript      string   `json:"transcript"`
	Locale          string   `json:"locale"`
	AcceptedLocales []string `json:"acceptedLocales"`
}

type ResearchResponse struct {
	Action          string   `json:"action"`
	Transcript      string   `json:"transcript"`
	Title           string   `json:"title"`
	Summary         string   `json:"summary"`
	KeyFindings     []string `json:"keyFindings"`
	Recommendation  string   `json:"recommendation"`
	Images          []string `json:"images"`
	SpokenAnswer    string   `json:"spokenAnswer"`
	FollowUpPrompts []string `json:"followUpPrompts"`
}

type BootstrapRequest struct {
	Platform   string `json:"platform"`
	Arch       string `json:"arch"`
	AppVersion string `json:"appVersion"`
	DeviceName string `json:"deviceName"`
}

type BootstrapResponse struct {
	User     UserProfile   `json:"user"`
	Runtime  RuntimeConfig `json:"runtime"`
	Features FeatureFlags  `json:"features"`
}

type UserProfile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type RuntimeConfig struct {
	APIBaseURL          string `json:"apiBaseUrl"`
	NativeWakeEnabled   bool   `json:"nativeWakeEnabled"`
	NativeWakeProvider  string `json:"nativeWakeProvider"`
	NativeWakeAccessKey string `json:"nativeWakeAccessKey"`
}

type FeatureFlags struct {
	ResearchEnabled    bool `json:"researchEnabled"`
	GoogleOAuthEnabled bool `json:"googleOAuthEnabled"`
}
