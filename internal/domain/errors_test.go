package domain

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorsIsBatchNotFound(t *testing.T) {
	err := ErrBatchNotFound
	// fmt.Println("errors.Is lot introuvable :", errors.Is(err, ErrBatchNotFound))

	if !errors.Is(err, ErrBatchNotFound) {
		t.Fatal("Erreur : errors.Is doit reconnaître ErrBatchNotFound")
	}
}

func TestErrorsAsValidationError(t *testing.T) {
	err := NewValidationError("urls", "liste vide")

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatal("Erreur : errors.As doit extraire *ValidationError")
	}
	if valErr.Field != "urls" {
		t.Fatalf("Erreur : champ attendu : urls. Champ obtenu : %s", valErr.Field)
	}
	// fmt.Println("errors.As champ : ", valErr.Field, valErr.Message)
}

func TestErrorsAsValidationErrorWrappee(t *testing.T) {
	// erreur wrappée comme on le fera depuis l'API ou le store
	err := fmt.Errorf("Erreur : création lot : %w", NewValidationError("parallelism", "doit être > 0"))

	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatal("Erreur : errors.As doit marcher sur une erreur wrappée avec %w")
	}
	if valErr.Field != "parallelism" {
		t.Fatalf("Erreur : champ attendu parallelism, obtenu %s", valErr.Field)
	}
}

func TestValidateBatchInput(t *testing.T) {
	err := ValidateBatchInput(nil, 2)
	// fmt.Println("validation urls vide :", err)
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("erreur validation attendue, obtenu %v", err)
	}

	err = ValidateBatchInput([]string{"https://exemple.fr"}, 0)
	if !errors.As(err, &valErr) || valErr.Field != "parallelism" {
		t.Fatalf("parallelism invalide attendu, obtenu %v", err)
	}

	err = ValidateBatchInput([]string{"https://exemple.fr"}, 3)
	if err != nil {
		t.Fatalf("entrée valide ne doit pas errer: %v", err)
	}
	// dernier cas du test
	// fmt.Println("ValidateBatchInput OK pour exemple.fr")
}
