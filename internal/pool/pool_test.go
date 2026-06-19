package pool

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Cengokill/Examen-final-go/internal/checker"
	"github.com/Cengokill/Examen-final-go/internal/domain"
)

func TestRunCollecteTousLesResultats(t *testing.T) {
	mock := &checker.MockChecker{Delay: 5 * time.Millisecond}
	runner := NewRunner(mock)

	urls := []string{
		"https://a.test",
		"https://b.test",
		"https://c.test",
	}

	results := runner.Run(context.Background(), urls, Options{
		Concurrency:  2,
		BatchTimeout: 2 * time.Second,
		URLTimeout:   time.Second,
	})

	if len(results) != len(urls) {
		t.Fatalf("%d résultats attendus, obtenu %d", len(urls), len(results))
	}
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

	// doit rester <= 2
	// fmt.Println("max goroutines actives :", mock.MaxActifs())
	if mock.MaxActifs() > concurrency {
		t.Fatalf("concurrence dépassée : max %d, limite %d", mock.MaxActifs(), concurrency)
	}
}

func TestRunBatchTimeout(t *testing.T) {
	mock := &checker.MockChecker{Delay: 200 * time.Millisecond}
	runner := NewRunner(mock)

	urls := []string{
		"https://lent-1.test",
		"https://lent-2.test",
		"https://lent-3.test",
		"https://lent-4.test",
	}

	results := runner.Run(context.Background(), urls, Options{
		Concurrency:  2,
		BatchTimeout: 50 * time.Millisecond,
		URLTimeout:   time.Second,
	})

	// toutes les URLs ne seront pas forcément traitées avant le timeout global
	if len(results) >= len(urls) {
		t.Fatalf("timeout global attendu, tous les résultats collectés (%d)", len(results))
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

	if len(results) != 1 {
		t.Fatalf("1 résultat attendu, obtenu %d", len(results))
	}
	if results[0].Error == "" {
		t.Fatal("erreur de timeout par URL attendue")
	}
}

func TestRunAnnulationContexte(t *testing.T) {
	mock := &checker.MockChecker{Delay: 100 * time.Millisecond}
	runner := NewRunner(mock)

	ctx, cancel := context.WithCancel(context.Background())

	urls := make([]string, 10)
	for i := range urls {
		urls[i] = fmt.Sprintf("https://cancel-%d.test", i)
	}

	done := make(chan []domain.CheckResult)
	go func() {
		done <- runner.Run(ctx, urls, Options{
			Concurrency:  2,
			BatchTimeout: 10 * time.Second,
			URLTimeout:   time.Second,
		})
	}()

	tempsDormir := 30 * time.Millisecond
	time.Sleep(tempsDormir)
	cancel()

	results := <-done
	if len(results) >= len(urls) {
		t.Fatalf("annulation attendue, %d/%d résultats", len(results), len(urls))
	}
}
