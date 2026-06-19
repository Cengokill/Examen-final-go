package api

import (
	"fmt"
	"strings"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// validateCreateRequest vérifie le contrat JSON de POST /v1/checks.
func validateCreateRequest(req createCheckRequest) (concurrency, timeoutMs int, err error) {
	if len(req.URLs) == 0 {
		// fmt.Println("validation : urls vide")
		return 0, 0, domain.NewValidationError("urls", "au moins une URL est requise")
	}
	if len(req.URLs) > maxURLs {
		return 0, 0, domain.NewValidationError("urls", fmt.Sprintf("maximum %d URLs", maxURLs))
	}

	for i, rawURL := range req.URLs {
		url := strings.TrimSpace(rawURL)
		if url == "" {
			return 0, 0, fmt.Errorf("urls[%d]: %w", i, domain.NewValidationError("urls", "URL vide"))
		}
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			// fmt.Println("url invalide index", i, " : ", rawURL)
			return 0, 0, fmt.Errorf("urls[%d]: %w", i, domain.NewValidationError("urls", "schéma http ou https requis"))
		}
	}

	concurrency = defaultConcurrency
	timeoutMs = defaultTimeoutMs
	if req.Options != nil {
		if req.Options.Concurrency != 0 {
			concurrency = req.Options.Concurrency
		}
		if req.Options.TimeoutMs != 0 {
			timeoutMs = req.Options.TimeoutMs
		}
	}

	if concurrency < minConcurrency || concurrency > maxConcurrency {
		return 0, 0, domain.NewValidationError("options.concurrency", fmt.Sprintf("doit être entre %d et %d", minConcurrency, maxConcurrency))
	}
	if timeoutMs < minTimeoutMs || timeoutMs > maxTimeoutMs {
		return 0, 0, domain.NewValidationError("options.timeout_ms", fmt.Sprintf("doit être entre %d et %d", minTimeoutMs, maxTimeoutMs))
	}

	// fmt.Println("validation OK :", len(req.URLs), "urls, concurrency", concurrency, "timeout_ms", timeoutMs)
	return concurrency, timeoutMs, nil
}
