package api

// createCheckRequest est le corps JSON de POST /v1/checks.
type createCheckRequest struct {
	URLs    []string            `json:"urls"`
	Options *createCheckOptions `json:"options"`
}

// createCheckOptions contient les options optionnelles du lot.
type createCheckOptions struct {
	Concurrency int `json:"concurrency"`
	TimeoutMs   int `json:"timeout_ms"`
}

// checkResponse est la réponse JSON d'un lot (POST 201 ou GET 200).
type checkResponse struct {
	BatchID   string           `json:"batch_id"`
	CreatedAt string           `json:"created_at"`
	Summary   summaryResponse  `json:"summary"`
	Results   []resultResponse `json:"results"`
}

// summaryResponse agrège les stats exposées par l'API.
type summaryResponse struct {
	Total      int   `json:"total"`
	Up         int   `json:"up"`
	Down       int   `json:"down"`
	DurationMs int64 `json:"duration_ms"`
}

// resultResponse est un résultat de vérification côté API.
type resultResponse struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code,omitempty"` // omitempty : absent si 0 (timeout DNS)
	OK         bool   `json:"ok"`
	LatencyMs  int64  `json:"latency_ms"`
	Error      string `json:"error,omitempty"`
}

// errorResponse est le format d'erreur uniforme de l'API.
type errorResponse struct {
	Error errorBody `json:"error"`
}

// errorBody décrit une erreur métier ou technique.
type errorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

const (
	defaultConcurrency = 8
	defaultTimeoutMs   = 5000
	// sans options dans le JSON
	maxURLs        = 100
	minConcurrency = 1
	maxConcurrency = 50
	minTimeoutMs   = 100
	maxTimeoutMs   = 30000
)
