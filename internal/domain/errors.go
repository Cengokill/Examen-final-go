package domain

import (
	"errors"
	"fmt"
)

// ErrBatchNotFound est renvoyée par Store.Get quand l'identifiant est inconnu.
var ErrBatchNotFound = errors.New("lot introuvable")

// ValidationError signale un champ invalide dans une requête (ex: urls, parallelism).
type ValidationError struct {
	Field   string
	Message string
}

// Error implémente l'interface error.
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation %s: %s", e.Field, e.Message)
}

// NewValidationError crée une erreur de validation sur un champ précis.
func NewValidationError(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
