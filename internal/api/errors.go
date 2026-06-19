package api

import (
	"errors"
	"net/http"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// StatusFromError traduit une erreur métier en code HTTP pour les handlers.
func StatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if errors.Is(err, domain.ErrBatchNotFound) {
		return http.StatusNotFound
	}

	var valErr *domain.ValidationError
	if errors.As(err, &valErr) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
