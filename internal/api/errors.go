package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// ErrorCode retourne le code d'erreur API à partir d'une erreur métier.
func ErrorCode(err error) string {
	if errors.Is(err, domain.ErrBatchNotFound) {
		return "batch_not_found"
	}

	var valErr *domain.ValidationError
	if errors.As(err, &valErr) {
		return "invalid_request"
	}

	return "internal"
}

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

// writeJSON encode une réponse JSON avec le status HTTP donné.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	// fmt.Println("writeJSON statut :", status)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// writeAPIError renvoie le corps d'erreur uniforme { "error": { "code", "message" } }.
func writeAPIError(w http.ResponseWriter, status int, code, message string) {
	// fmt.Println("API error", statut, code, " : ", message)
	writeJSON(w, status, errorResponse{
		Error: errorBody{
			Code:    code,
			Message: message,
		},
	})
}

// writeDomainError traduit une erreur domaine en réponse JSON.
func writeDomainError(w http.ResponseWriter, err error) {
	status := StatusFromError(err)
	code := ErrorCode(err)
	message := err.Error()

	if errors.Is(err, domain.ErrBatchNotFound) {
		message = "aucun lot avec l'id demandé"
	}
	// test errors.Is / errors.As
	// fmt.Println("writeDomainError :", code, status, message)

	writeAPIError(w, status, code, message)
}
