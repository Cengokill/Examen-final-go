package checker

import (
	"context"
	"sync"
	"time"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// MockChecker simule des vérifications déterministes pour les tests du pool.
type MockChecker struct {
	Delay    time.Duration
	Response func(url string) domain.CheckResult

	mu        sync.Mutex
	actifs    int
	maxActifs int
}

// Check simule une requête et suit le nombre de goroutines actives (tests concurrence).
func (m *MockChecker) Check(ctx context.Context, url string) domain.CheckResult {
	m.mu.Lock()
	m.actifs++
	if m.actifs > m.maxActifs {
		m.maxActifs = m.actifs
	}
	// fmt.Println("mock actifs :", m.actifs, "max :", m.maxActifs) // debug pool concurrency=2
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		m.actifs--
		m.mu.Unlock()
	}()

	if m.Delay > 0 {
		timer := time.NewTimer(m.Delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			// fmt.Println("mock annulé :", url, ctx.Err()) // test batch timeout pool
			return domain.CheckResult{
				URL:       url,
				Available: false,
				Error:     ctx.Err().Error(),
			}
		case <-timer.C:
		}
	}

	if m.Response != nil {
		return m.Response(url)
	}

	// fmt.Println("mock OK par défaut : ", url, "delay", m.Delay)
	return domain.CheckResult{
		URL:       url,
		Available: true,
		LatencyMs: m.Delay.Milliseconds(),
	}
}

// MaxActifs retourne le pic de goroutines simultanées observé pendant les checks.
func (m *MockChecker) MaxActifs() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	// fmt.Println("mock MaxActifs : ", m.maxActifs)
	return m.maxActifs
}
