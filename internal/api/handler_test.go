package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Cengokill/Examen-final-go/internal/checker"
	"github.com/Cengokill/Examen-final-go/internal/domain"
	"github.com/Cengokill/Examen-final-go/internal/pool"
)

// mockStore implémente domain.Store pour httptest (sans BDD réelle).
type mockStore struct {
	batches map[string]domain.Batch
}

func newMockStore() *mockStore {
	return &mockStore{batches: make(map[string]domain.Batch)}
}

func (m *mockStore) Save(ctx context.Context, b domain.Batch) error {
	_ = ctx
	m.batches[b.ID] = b
	return nil
}

func (m *mockStore) Get(ctx context.Context, id string) (domain.Batch, error) {
	_ = ctx
	b, ok := m.batches[id]
	if !ok {
		return domain.Batch{}, domain.ErrBatchNotFound
	}
	return b, nil
}

func newTestServerWith(store domain.Store) *Server {
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
	return NewServer(store, pool.NewRunner(mock))
}

func newTestServer() *Server {
	return newTestServerWith(newMockStore())
}

func TestHandlersHTTPTable(t *testing.T) {
	store := newMockStore()
	srv := newTestServerWith(store)
	logger := NewJSONLogger("ERROR")
	handler := srv.Handler(logger)

	// POST réussi pour alimenter le GET suivant
	postBody := `{"urls":["https://go.dev"],"options":{"concurrency":2,"timeout_ms":1000}}`
	postReq := httptest.NewRequest(http.MethodPost, "/v1/checks", bytes.NewBufferString(postBody))
	postRec := httptest.NewRecorder()
	handler.ServeHTTP(postRec, postReq)

	if postRec.Code != http.StatusCreated {
		t.Fatalf("POST status attendu 201, obtenu %d", postRec.Code)
	}

	var created checkResponse
	if err := json.NewDecoder(postRec.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}

	cas := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantCode   string
	}{
		{
			name:       "POST /v1/checks réussi",
			method:     http.MethodPost,
			path:       "/v1/checks",
			body:       `{"urls":["https://killian.test"]}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "GET /v1/checks/{id} réussi",
			method:     http.MethodGet,
			path:       "/v1/checks/" + created.BatchID,
			wantStatus: http.StatusOK,
		},
		{
			name:       "GET lot inconnu 404",
			method:     http.MethodGet,
			path:       "/v1/checks/b_inconnu",
			wantStatus: http.StatusNotFound,
			wantCode:   "batch_not_found",
		},
		{
			name:       "POST urls vide 400",
			method:     http.MethodPost,
			path:       "/v1/checks",
			body:       `{"urls":[]}`,
			wantStatus: http.StatusBadRequest,
			wantCode:   "invalid_request",
		},
	}

	for _, tc := range cas {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest(tc.method, tc.path, bytes.NewBufferString(tc.body))
			} else {
				req = httptest.NewRequest(tc.method, tc.path, nil)
			}
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status attendu %d, obtenu %d body=%s", tc.wantStatus, rec.Code, rec.Body.String())
			}

			if tc.wantCode != "" {
				var errResp errorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
					t.Fatal(err)
				}
				if errResp.Error.Code != tc.wantCode {
					t.Fatalf("code attendu %s, obtenu %s", tc.wantCode, errResp.Error.Code)
				}
			}
		})
	}
	// fmt.Println("httptest table : POST 201 + GET 404 OK")
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
}
