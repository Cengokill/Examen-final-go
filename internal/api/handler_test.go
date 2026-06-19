package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Cengokill/Examen-final-go/internal/checker"
	"github.com/Cengokill/Examen-final-go/internal/domain"
	"github.com/Cengokill/Examen-final-go/internal/pool"
	"github.com/Cengokill/Examen-final-go/internal/store"
)

func newTestServer() *Server {
	mock := &checker.MockChecker{
		Delay: 2 * time.Millisecond,
		Response: func(url string) domain.CheckResult {
			return domain.CheckResult{
				URL:        url,
				StatusCode: 200,
				Available:  true,
				LatencyMs:  10,
			}
		},
	}
	return NewServer(store.NewMemoryStore(), pool.NewRunner(mock))
}

func TestHealthz(t *testing.T) {
	srv := newTestServer()
	logger := NewJSONLogger("ERROR")
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	srv.Handler(logger).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status attendu 200, obtenu %d", rec.Code)
	}
	// fmt.Println("healthz body :", rec.Body.String()) // {"status":"ok"}
}

func TestCreateCheck(t *testing.T) {
	srv := newTestServer()
	logger := NewJSONLogger("ERROR")

	body := `{"urls":["https://go.dev"],"options":{"concurrency":2,"timeout_ms":1000}}`
	req := httptest.NewRequest(http.MethodPost, "/v1/checks", bytes.NewBufferString(body))
	rec := httptest.NewRecorder()

	srv.Handler(logger).ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status attendu 201, obtenu %d body=%s", rec.Code, rec.Body.String())
	}

	var resp checkResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp.BatchID == "" {
		t.Fatal("batch_id manquant")
	}
	if resp.Summary.Total != 1 {
		t.Fatalf("summary.total attendu 1, obtenu %d", resp.Summary.Total)
	}
	// fmt.Println("POST /v1/checks :", resp.BatchID, resp.Summary) // test curl
}

func TestGetCheckNotFound(t *testing.T) {
	srv := newTestServer()
	logger := NewJSONLogger("ERROR")

	req := httptest.NewRequest(http.MethodGet, "/v1/checks/b_inconnu", nil)
	rec := httptest.NewRecorder()

	srv.Handler(logger).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status attendu 404, obtenu %d", rec.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatal(err)
	}
	if errResp.Error.Code != "batch_not_found" {
		t.Fatalf("code attendu batch_not_found, obtenu %s", errResp.Error.Code)
	}
	// fmt.Println("404 error body :", errResp.Error.Message)
}

func TestRecoveryMiddleware(t *testing.T) {
	logger := NewJSONLogger("ERROR")
	handler := recoveryMiddleware(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status attendu 500, obtenu %d", rec.Code)
	}

	var errResp errorResponse
	if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
		t.Fatal(err)
	}
	if errResp.Error.Code != "internal" {
		t.Fatalf("code attendu internal, obtenu %s", errResp.Error.Code)
	}
	// fmt.Println("recovery 500 :", errResp.Error.Code)
}
