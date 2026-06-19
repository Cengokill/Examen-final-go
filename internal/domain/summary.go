package domain

import "time"

// ComputeSummary calcule le résumé à partir d'une slice de CheckResult.
// On utilise une map pour compter les dispo / échecs de façon idiomatique.
func ComputeSummary(results []CheckResult) BatchSummary {
	compteurs := map[string]int{
		"available": 0,
		"failed":    0,
	}

	var dureeTotale int64
	for _, r := range results {
		if r.Available {
			compteurs["available"]++
		} else {
			compteurs["failed"]++
		}
		dureeTotale += r.LatencyMs
	}

	return BatchSummary{
		Total:           len(results),
		Available:       compteurs["available"],
		Failed:          compteurs["failed"],
		TotalDurationMs: dureeTotale,
	}
}

// NewBatch construit un lot avec la date de création et le résumé calculé.
func NewBatch(id string, results []CheckResult) Batch {
	// copie la slice pour éviter les surprises si le caller la modifie après
	copie := make([]CheckResult, len(results))
	copy(copie, results)

	return Batch{
		ID:        id,
		CreatedAt: time.Now().UTC(),
		Results:   copie,
		Summary:   ComputeSummary(copie),
	}
}
