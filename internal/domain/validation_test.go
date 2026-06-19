package domain

import (
	"errors"
	"testing"
)

func TestValidateBatchInputTable(t *testing.T) {
	// pas de "url" au lieu de "urls"
	cas := []struct {
		name        string
		urls        []string
		parallelism int
		wantField   string
		wantErr     bool
	}{
		{
			name:        "urls vide",
			urls:        nil,
			parallelism: 2,
			wantField:   "urls",
			wantErr:     true,
		},
		{
			name:        "url sans schéma",
			urls:        []string{"ftp://bad.test"},
			parallelism: 2,
			wantField:   "urls",
			wantErr:     true,
		},
		{
			name:        "parallelism nul",
			urls:        []string{"https://ok.test"},
			parallelism: 0,
			wantField:   "parallelism",
			wantErr:     true,
		},
		{
			name:        "entrée valide",
			urls:        []string{"https://go.dev", "http://localhost:8080/healthz"},
			parallelism: 4,
			wantErr:     false,
		},
	}

	for _, tc := range cas {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateBatchInput(tc.urls, tc.parallelism)
			if tc.wantErr {
				if err == nil {
					t.Fatal("erreur attendue")
				}
				var valErr *ValidationError
				if !errors.As(err, &valErr) {
					// bloqué ici au début : fallait errors.As pas errors.Is pour ValidationError
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
		})
	}
	// fmt.Println("ValidateBatchInput table : 4 cas OK")
}
