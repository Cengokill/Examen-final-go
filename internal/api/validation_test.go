package api

import "testing"

func TestValidateCreateRequest(t *testing.T) {
	_, _, err := validateCreateRequest(createCheckRequest{})
	if err == nil {
		t.Fatal("urls vide doit échouer")
	}
	// fmt.Println("urls vide : ", err)

	tooMany := make([]string, 101)
	for i := range tooMany {
		tooMany[i] = "https://exemple.test"
	}
	_, _, err = validateCreateRequest(createCheckRequest{URLs: tooMany})
	if err == nil {
		t.Fatal("plus de 100 urls doit échouer")
	}
	// fmt.Println("101 urls rejetées : ", err) // max 100

	concurrency, timeoutMs, err := validateCreateRequest(createCheckRequest{
		URLs: []string{"https://go.dev"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if concurrency != defaultConcurrency || timeoutMs != defaultTimeoutMs {
		t.Fatalf("défauts attendus %d/%d, obtenu %d/%d", defaultConcurrency, defaultTimeoutMs, concurrency, timeoutMs)
	}
	// fmt.Println("défauts validation : ", concurrency, timeoutMs)
}
