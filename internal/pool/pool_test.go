package pool

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/Cengokill/Examen-final-go/internal/checker"
	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// deterministicMock renvoie ok/ko selon l'URL (sans réseau).
func deterministicMock() *checker.MockChecker {
	return &checker.MockChecker{
		Delay: 5 * time.Millisecond,
		Response: func(url string) domain.CheckResult {
			if strings.Contains(url, "fail") {
				return domain.CheckResult{
					URL:       url,
					Available: false,
					Error:     "mock: host down",
					LatencyMs: 5,
				}
			}
			return domain.CheckResult{
				URL:        url,
				StatusCode: 200,
				Available:  true,
				LatencyMs:  5,
			}
		},
	}
}

func TestRunTable(t *testing.T) {
	cas := []struct {
		name         string
		urls         []string
		concurrency  int
		batchTimeout time.Duration
		urlTimeout   time.Duration
		wantCount    int
		wantUp       int
		wantDown     int
	}{
		{
			name:         "collecte complète",
			urls:         []string{"https://a.test", "https://b.test", "https://c.test"},
			concurrency:  2,
			batchTimeout: 2 * time.Second,
			urlTimeout:   time.Second,
			wantCount:    3,
			wantUp:       3,
		},
		{
			name:         "mock déterministe ok/ko",
			urls:         []string{"https://ok.test", "https://fail.test"},
			concurrency:  1,
			batchTimeout: 2 * time.Second,
			urlTimeout:   time.Second,
			wantCount:    2,
			wantUp:       1,
			wantDown:     1,
		},
	}

	for _, tc := range cas {
		t.Run(tc.name, func(t *testing.T) {
			runner := NewRunner(deterministicMock())
			results := runner.Run(context.Background(), tc.urls, Options{
				Concurrency:  tc.concurrency,
				BatchTimeout: tc.batchTimeout,
				URLTimeout:   tc.urlTimeout,
			})

			if len(results) != tc.wantCount {
				t.Fatalf("%d résultats attendus, obtenu %d", tc.wantCount, len(results))
			}

			summary := domain.ComputeSummary(results)
			if summary.Available != tc.wantUp {
				t.Fatalf("up attendu %d, obtenu %d", tc.wantUp, summary.Available)
			}
			if summary.Failed != tc.wantDown {
				t.Fatalf("down attendu %d, obtenu %d", tc.wantDown, summary.Failed)
			}
		})
	}
	// fmt.Println("pool table : mock déterministe sans réseau") // fan-out/fan-in
}

func TestRunConcurrencyBornee(t *testing.T) {
	mock := &checker.MockChecker{Delay: 30 * time.Millisecond}
	runner := NewRunner(mock)

	const concurrency = 2
	urls := make([]string, 8)
	for i := range urls {
		urls[i] = fmt.Sprintf("https://site-%d.test", i)
	}

	runner.Run(context.Background(), urls, Options{
		Concurrency:  concurrency,
		BatchTimeout: 5 * time.Second,
		URLTimeout:   time.Second,
	})

	if mock.MaxActifs() > concurrency {
		t.Fatalf("concurrence dépassée : max %d, limite %d", mock.MaxActifs(), concurrency)
	}
}

func TestRunContextAnnulationOuTimeout(t *testing.T) {
	cas := []struct {
		name       string
		cancel     bool
		batchTO    time.Duration
		wantPartiel bool
	}{
		{
			name:        "annulation manuelle",
			cancel:      true,
			batchTO:     10 * time.Second,
			wantPartiel: true,
		},
		{
			name:        "timeout global lot",
			cancel:      false,
			batchTO:     50 * time.Millisecond,
			wantPartiel: true,
		},
	}

	for _, tc := range cas {
		t.Run(tc.name, func(t *testing.T) {
			mock := &checker.MockChecker{Delay: 100 * time.Millisecond}
			runner := NewRunner(mock)

			ctx, cancel := context.WithCancel(context.Background())
			urls := make([]string, 10)
			for i := range urls {
				urls[i] = fmt.Sprintf("https://lent-%d.test", i)
			}

			done := make(chan []domain.CheckResult)
			go func() {
				done <- runner.Run(ctx, urls, Options{
					Concurrency:  2,
					BatchTimeout: tc.batchTO,
					URLTimeout:   time.Second,
				})
			}()

			if tc.cancel {
				time.Sleep(30 * time.Millisecond)
				cancel()
			}

			results := <-done
			if tc.wantPartiel && len(results) >= len(urls) {
				t.Fatalf("résultats partiels attendus, obtenu %d/%d", len(results), len(urls))
			}
		})
	}
}

func TestRunURLTimeout(t *testing.T) {
	mock := &checker.MockChecker{
		Delay: 150 * time.Millisecond,
		Response: func(url string) domain.CheckResult {
			return domain.CheckResult{URL: url, Available: true}
		},
	}
	runner := NewRunner(mock)

	results := runner.Run(context.Background(), []string{"https://timeout.test"}, Options{
		Concurrency:  1,
		BatchTimeout: 2 * time.Second,
		URLTimeout:   20 * time.Millisecond,
	})

	if len(results) != 1 || results[0].Error == "" {
		t.Fatal("erreur de timeout par URL attendue")
	}
}
