package domain

import "testing"

func TestComputeSummary(t *testing.T) {
	results := []CheckResult{
		{URL: "https://ok.test", StatusCode: 200, Available: true, LatencyMs: 100},
		{URL: "https://ok2.test", StatusCode: 201, Available: true, LatencyMs: 50},
		{URL: "https://ko.test", StatusCode: 0, Available: false, LatencyMs: 200, Error: "timeout"},
	}

	summary := ComputeSummary(results)

	// fmt.Println("summary test : ", summary.Total, summary.Available, summary.Failed)
	if summary.Total != 3 {
		t.Fatalf("Total attendu 3, obtenu %d", summary.Total)
	}
	if summary.Available != 2 {
		t.Fatalf("Available attendu 2, obtenu %d", summary.Available)
	}
	if summary.Failed != 1 {
		t.Fatalf("Failed attendu 1, obtenu %d", summary.Failed)
	}
	if summary.TotalDurationMs != 350 {
		t.Fatalf("TotalDurationMs attendu 350, obtenu %d", summary.TotalDurationMs)
	}
}

func TestComputeSummaryVide(t *testing.T) {
	summary := ComputeSummary(nil)

	if summary.Total != 0 || summary.Available != 0 || summary.Failed != 0 {
		t.Fatalf("résumé vide attendu, obtenu : %+v", summary)
	}
}

func TestNewBatch(t *testing.T) {
	results := []CheckResult{
		{URL: "https://exemple.fr", Available: true, LatencyMs: 42},
	}

	batch := NewBatch("batch-1", results)

	// fmt.Println("NewBatch test : ", batch.ID, batch.Summary.Total)
	if batch.ID != "batch-1" {
		t.Fatalf("ID attendu batch-1, obtenu : %s", batch.ID)
	}
	if batch.CreatedAt.IsZero() {
		t.Fatal("CreatedAt ne doit pas être zéro")
	}
	if len(batch.Results) != 1 {
		t.Fatalf("1 résultat attendu, obtenu %d", len(batch.Results))
	}
	if batch.Summary.Total != 1 {
		t.Fatalf("Summary.Total attendu 1, obtenu %d", batch.Summary.Total)
	}

	// la slice d'origine ne doit pas impacter le batch
	results[0].URL = "modifié"
	if batch.Results[0].URL != "https://exemple.fr" {
		t.Fatal("le batch doit avoir sa propre copie des résultats")
	}
}
