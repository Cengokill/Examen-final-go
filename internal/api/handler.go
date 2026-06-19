package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Cengokill/Examen-final-go/internal/domain"
	"github.com/Cengokill/Examen-final-go/internal/pool"
)

// Server expose les handlers HTTP de URLWatch.
type Server struct {
	store  domain.Store
	runner *pool.Runner
}

// NewServer assemble le store et le runner.
func NewServer(store domain.Store, runner *pool.Runner) *Server {
	return &Server{
		store:  store,
		runner: runner,
	}
}

// Handler retourne le routeur HTTP avec middlewares (recovery + logging).
func (s *Server) Handler(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.handleHealthz)
	mux.HandleFunc("/v1/checks", s.handleChecks)
	mux.HandleFunc("/v1/checks/", s.handleCheckByID)

	return Chain(
		recoveryMiddleware(logger),
		loggingMiddleware(logger),
	)(mux)
}

// handleHealthz répond à la sonde de vivacité.
func (s *Server) handleHealthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAPIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "méthode non autorisée")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleChecks gère POST /v1/checks.
func (s *Server) handleChecks(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/v1/checks" {
		writeAPIError(w, http.StatusNotFound, "invalid_request", "chemin inconnu")
		return
	}

	switch r.Method {
	case http.MethodPost:
		s.createCheck(w, r)
	default:
		writeAPIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "méthode non autorisée")
	}
}

// handleCheckByID gère GET /v1/checks/{id}.
func (s *Server) handleCheckByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeAPIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "méthode non autorisée")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/v1/checks/")
	if id == "" || strings.Contains(id, "/") {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "identifiant de lot manquant")
		return
	}

	batch, err := s.store.Get(r.Context(), id)
	if err != nil {
		// fmt.Println("GET lot inconnu : ", id) // curl b_inconnu : 404
		writeAPIError(w, http.StatusNotFound, "batch_not_found", "aucun lot avec l'id "+id)
		return
	}

	writeJSON(w, http.StatusOK, toCheckResponse(batch))
}

// createCheck exécute un lot, le persiste et renvoie 201 Created.
func (s *Server) createCheck(w http.ResponseWriter, r *http.Request) {
	var req createCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeAPIError(w, http.StatusBadRequest, "invalid_request", "JSON invalide")
		return
	}

	concurrency, timeoutMs, err := validateCreateRequest(req)
	if err != nil {
		// fmt.Println("validation POST échouée :", err) // test urls:[] ou concurrency:99
		writeDomainError(w, err)
		return
	}

	urls := make([]string, len(req.URLs))
	for i, raw := range req.URLs {
		urls[i] = strings.TrimSpace(raw)
	}

	batchID := newBatchID()
	urlTimeout := time.Duration(timeoutMs) * time.Millisecond
	batchTimeout := computeBatchTimeout(len(urls), concurrency, timeoutMs)
	// fmt.Println("pool run :", len(urls), "urls, concurrency", concurrency, "timeout", timeoutMs, "ms")

	results := s.runner.Run(r.Context(), urls, pool.Options{
		Concurrency:  concurrency,
		BatchTimeout: batchTimeout,
		URLTimeout:   urlTimeout,
	})

	batch := domain.NewBatch(batchID, results)
	if err := s.store.Save(r.Context(), batch); err != nil {
		writeAPIError(w, http.StatusInternalServerError, "internal", "impossible de persister le lot")
		return
	}

	// fmt.Println("lot créé :", batchID, "urls :", len(req.URLs)) // test curl POST
	// fmt.Println("résumé :", batch.Summary.Available, "up /", batch.Summary.Failed, "down")
	if rec, ok := w.(*responseRecorder); ok {
		rec.batchID = batchID
	}

	writeJSON(w, http.StatusCreated, toCheckResponse(batch))
}

// newBatchID génère un identifiant du type b_4f3c1a.
func newBatchID() string {
	buf := make([]byte, 3)
	_, _ = rand.Read(buf)
	return "b_" + hex.EncodeToString(buf)
}

// computeBatchTimeout estime le délai global du lot (non exposé dans l'API).
func computeBatchTimeout(urlCount, concurrency, timeoutMs int) time.Duration {
	waves := (urlCount + concurrency - 1) / concurrency
	ms := int64(timeoutMs) * int64(waves+1)
	return time.Duration(ms) * time.Millisecond
}

// toCheckResponse mappe un Batch domaine vers le contrat JSON API.
func toCheckResponse(batch domain.Batch) checkResponse {
	results := make([]resultResponse, len(batch.Results))
	for i, r := range batch.Results {
		results[i] = resultResponse{
			URL:        r.URL,
			StatusCode: r.StatusCode,
			OK:         r.Available,
			LatencyMs:  r.LatencyMs,
			Error:      r.Error,
		}
	}

	return checkResponse{
		BatchID:   batch.ID,
		CreatedAt: batch.CreatedAt.UTC().Format(time.RFC3339),
		Summary: summaryResponse{
			Total:      batch.Summary.Total,
			Up:         batch.Summary.Available,
			Down:       batch.Summary.Failed,
			DurationMs: batch.Summary.TotalDurationMs,
		},
		Results: results,
	}
}
