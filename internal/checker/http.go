package checker

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// HTTPChecker vérifie une URL via un vrai appel HTTP GET.
type HTTPChecker struct {
	client *http.Client
}

// NewHTTPChecker crée un checker HTTP avec le client par défaut.
func NewHTTPChecker() *HTTPChecker {
	return &HTTPChecker{
		client: &http.Client{},
	}
}

// Check exécute la requête en respectant le context (timeout / annulation).
func (h *HTTPChecker) Check(ctx context.Context, url string) domain.CheckResult {
	debut := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return domain.CheckResult{
			URL:       url,
			Available: false,
			LatencyMs: time.Since(debut).Milliseconds(),
			Error:     err.Error(),
		}
	}

	resp, err := h.client.Do(req)
	latence := time.Since(debut).Milliseconds()
	if err != nil {
		// fmt.Println("erreur HTTP", url, "→", err) // test avec timeout 1ms sur url lente
		return domain.CheckResult{
			URL:       url,
			Available: false,
			LatencyMs: latence,
			Error:     err.Error(),
		}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	disponible := resp.StatusCode >= 200 && resp.StatusCode < 400
	// fmt.Println(url, "→", resp.StatusCode, disponible, latence, "ms") // test manuel https://httpbin.org/get
	return domain.CheckResult{
		URL:        url,
		StatusCode: resp.StatusCode,
		Available:  disponible,
		LatencyMs:  latence,
	}
}
