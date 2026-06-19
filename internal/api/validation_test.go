package api

import (
	"errors"
	"testing"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

func TestValidateCreateRequestTable(t *testing.T) {
	cas := []struct {
		name            string
		req             createCheckRequest
		wantErr         bool
		wantField       string
		wantConcurrency int
		wantTimeoutMs   int
	}{
		{
			name:      "urls vide",
			req:       createCheckRequest{},
			wantErr:   true,
			wantField: "urls",
		},
		{
			name: "trop d'urls",
			req: createCheckRequest{
				URLs: make([]string, 101),
			},
			wantErr:   true,
			wantField: "urls",
		},
		{
			name: "concurrency hors bornes",
			req: createCheckRequest{
				URLs:    []string{"https://go.dev"},
				Options: &createCheckOptions{Concurrency: 99},
			},
			wantErr:   true,
			wantField: "options.concurrency",
		},
		{
			name: "timeout_ms hors bornes",
			req: createCheckRequest{
				URLs:    []string{"https://go.dev"},
				Options: &createCheckOptions{TimeoutMs: 50},
			},
			wantErr:   true,
			wantField: "options.timeout_ms",
		},
		{
			name: "défauts sans options",
			req: createCheckRequest{
				URLs: []string{"https://go.dev"},
			},
			wantConcurrency: defaultConcurrency,
			wantTimeoutMs:   defaultTimeoutMs,
		},
		{
			name: "options explicites",
			req: createCheckRequest{
				URLs: []string{"https://go.dev"},
				Options: &createCheckOptions{
					Concurrency: 4,
					TimeoutMs:   2000,
				},
			},
			wantConcurrency: 4,
			wantTimeoutMs:   2000,
		},
	}

	for _, tc := range cas {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "trop d'urls" {
				for i := range tc.req.URLs {
					tc.req.URLs[i] = "https://killian.test"
				}
			}

			concurrency, timeoutMs, err := validateCreateRequest(tc.req)
			if tc.wantErr {
				if err == nil {
					t.Fatal("erreur attendue")
				}
				var valErr *domain.ValidationError
				if !errors.As(err, &valErr) {
					t.Fatalf("ValidationError attendue, obtenu %v", err)
				}
				if valErr.Field != tc.wantField {
					t.Fatalf("champ attendu %s, obtenu %s", tc.wantField, valErr.Field)
				}
				return
			}
			if err != nil {
				t.Fatalf("entrée valide ne doit pas errer : %v", err)
			}
			if concurrency != tc.wantConcurrency || timeoutMs != tc.wantTimeoutMs {
				t.Fatalf("options attendues %d/%d, obtenu %d/%d",
					tc.wantConcurrency, tc.wantTimeoutMs, concurrency, timeoutMs)
			}
		})
	}
}
