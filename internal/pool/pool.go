package pool

import (
	"context"
	"sync"
	"time"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// Options configure le worker pool (concurrence et timeouts).
type Options struct {
	Concurrency  int
	BatchTimeout time.Duration
	URLTimeout   time.Duration
}

// Runner orchestre la vérification concurrente d'un lot d'URLs.
type Runner struct {
	checker domain.Checker
}

// NewRunner crée un runner branché sur un Checker (HTTP ou mock).
func NewRunner(checker domain.Checker) *Runner {
	return &Runner{checker: checker}
}

// Run vérifie toutes les URLs avec un pool borné (fan-out / fan-in).
func (r *Runner) Run(ctx context.Context, urls []string, opts Options) []domain.CheckResult {
	if len(urls) == 0 {
		return nil
	}

	batchCtx, cancel := context.WithTimeout(ctx, opts.BatchTimeout)

	// Ajout d'un defer pour être sûr qu'onc ancel
	defer cancel()

	// jobs bufferisé : le fan-out peut envoyer sans attendre un worker libre
	jobs := make(chan string, len(urls))
	// results bufferisé : évite le deadlock si plusieurs workers finissent en même temps (TP 4a)
	results := make(chan domain.CheckResult, len(urls))

	var wg sync.WaitGroup

	// fan-out : distribution des URLs vers les workers
	go distribuerURLs(batchCtx, urls, jobs)

	// démarrage des workers (nombre fixe = concurrency, jamais une goroutine par URL)
	// fmt.Println("pool : lancement de", opts.Concurrency, "workers pour", len(urls), "urls")
	for w := 1; w <= opts.Concurrency; w++ {
		wg.Add(1)
		go worker(batchCtx, r.checker, opts.URLTimeout, jobs, results, &wg)
	}

	// fan-in : fermeture de results quand tous les workers ont terminé
	go func() {
		wg.Wait()
		close(results)
	}()

	collected := make([]domain.CheckResult, 0, len(urls))
	for result := range results {
		collected = append(collected, result)
	}

	// fmt.Println("pool terminé :", len(collected), "résultats sur", len(urls), "urls")
	return collected
}

// distribuerURLs envoie les URLs sur jobs puis ferme le canal (fan-out).
func distribuerURLs(ctx context.Context, urls []string, jobs chan<- string) {
	defer close(jobs)

	for _, url := range urls {
		select {
		case <-ctx.Done():
			// fmt.Println("fan-out interrompu :", ctx.Err())
			return
		case jobs <- url:
		}
	}
}

// worker qui lit les jobs et appelle le checker avec un timeout par URL.
func worker(
	ctx context.Context,
	checker domain.Checker,
	urlTimeout time.Duration,
	jobs <-chan string,
	results chan<- domain.CheckResult,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for url := range jobs {
		select {
		case <-ctx.Done():
			// fmt.Println("worker arrêté :", ctx.Err())
			return
		default:
		}

		urlCtx, cancel := context.WithTimeout(ctx, urlTimeout)
		result := checker.Check(urlCtx, url)
		cancel()

		select {
		case <-ctx.Done():
			return
		case results <- result:
		}
	}
}
