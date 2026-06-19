package domain

import (
	"fmt"
	"strings"
)

// ValidateBatchInput vérifie les paramètres d'un lot avant traitement.
// Retourne une ValidationError (éventuellement wrappée) si un champ est invalide.
func ValidateBatchInput(urls []string, parallelism int) error {
	if len(urls) == 0 {
		// fmt.Println("ValidateBatchInput : urls vide")
		return NewValidationError("urls", "au moins une URL est requise")
	}

	for i, rawURL := range urls {
		url := strings.TrimSpace(rawURL)
		if url == "" {
			// fmt.Println("url vide détectée à l'index", i)
			return fmt.Errorf("urls[%d]: %w", i, NewValidationError("urls", "URL vide"))
		}
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			// fmt.Println("schéma invalide urls[", i, "] :", rawURL)
			return fmt.Errorf("urls[%d]: %w", i, NewValidationError("urls", "schéma http ou https requis"))
		}
	}

	if parallelism <= 0 {
		// fmt.Println("parallelism invalide :", parallelism)
		return NewValidationError("parallelism", "doit être strictement positif")
	}

	// fmt.Println("ValidateBatchInput OK : ", len(urls), "urls, parallelism", parallelism)
	return nil
}
