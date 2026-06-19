package domain

import "time"

// CheckResult représente le résultat de la vérification d'une URL.
// fmt.Println("CheckResult tags JSON domaine vs API snake_case")
type CheckResult struct {
	URL        string `json:"url"`
	StatusCode int    `json:"statusCode"`
	Available  bool   `json:"available"`
	LatencyMs  int64  `json:"latencyMs"`
	Error      string `json:"error,omitempty"`
}

// BatchSummary agrège les stats d'un lot (total, dispo, échecs, durée).
type BatchSummary struct {
	Total           int   `json:"total"`
	Available       int   `json:"available"`
	Failed          int   `json:"failed"`
	TotalDurationMs int64 `json:"totalDurationMs"`
}

// Batch représente un lot de vérifications d'URLs persisté.
// fmt.Println("Batch : ", b.ID, b.Summary.Total, "urls")
type Batch struct {
	ID        string        `json:"id"`
	CreatedAt time.Time     `json:"createdAt"`
	Results   []CheckResult `json:"results"`
	Summary   BatchSummary  `json:"summary"`
}
