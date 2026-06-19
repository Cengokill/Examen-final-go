package api

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

func TestStatusFromError(t *testing.T) {
	casDerreur := []struct {
		name   string
		err    error
		status int
	}{
		{
			name:   "lot non trouvé",
			err:    domain.ErrBatchNotFound,
			status: 404,
		},
		{
			name:   "validation directe",
			err:    domain.NewValidationError("urls", "vide"),
			status: 400,
		},
		{
			name:   "validation wrappée",
			err:    fmt.Errorf("handler POST /batches: %w", domain.NewValidationError("parallelism", "invalide")),
			status: 400,
		},
		{
			name:   "erreur autre inconnue",
			err:    errors.New("panic sqlite"),
			status: 500,
		},
	}

	for _, tc := range casDerreur {
		t.Run(tc.name, func(t *testing.T) {
			got := StatusFromError(tc.err)
			if got != tc.status {
				t.Fatalf("statut attendu : %d, Mais obtenu : %d", tc.status, got)
			}
		})
	}
}
