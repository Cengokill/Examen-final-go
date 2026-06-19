package domain

import "testing"

func TestComputeSummaryTable(t *testing.T) {
	cas := []struct {
		name     string
		results  []CheckResult
		wantTotal int
		wantUp   int
		wantDown int
		wantMs   int64
	}{
		{
			name: "mixte up/down",
			results: []CheckResult{
				{Available: true, LatencyMs: 100},
				{Available: true, LatencyMs: 50},
				{Available: false, LatencyMs: 200},
			},
			wantTotal: 3,
			wantUp:    2,
			wantDown:  1,
			wantMs:    350,
		},
		{
			name:      "slice vide",
			results:   nil,
			wantTotal: 0,
			wantUp:    0,
			wantDown:  0,
			wantMs:    0,
		},
		{
			name: "tout en échec",
			results: []CheckResult{
				{Available: false, LatencyMs: 10},
				{Available: false, LatencyMs: 20},
			},
			wantTotal: 2,
			wantUp:    0,
			wantDown:  2,
			wantMs:    30,
		},
	}

	for _, tc := range cas {
		t.Run(tc.name, func(t *testing.T) {
			summary := ComputeSummary(tc.results)
			if summary.Total != tc.wantTotal {
				t.Fatalf("Total attendu %d, obtenu %d", tc.wantTotal, summary.Total)
			}
			if summary.Available != tc.wantUp {
				t.Fatalf("Available attendu %d, obtenu %d", tc.wantUp, summary.Available)
			}
			if summary.Failed != tc.wantDown {
				t.Fatalf("Failed attendu %d, obtenu %d", tc.wantDown, summary.Failed)
			}
			if summary.TotalDurationMs != tc.wantMs {
				t.Fatalf("TotalDurationMs attendu %d, obtenu %d", tc.wantMs, summary.TotalDurationMs)
			}
		})
	}
}

func TestNewBatch(t *testing.T) {
	results := []CheckResult{
		{URL: "https://exemple.fr", Available: true, LatencyMs: 42},
	}

	batch := NewBatch("batch-1", results)

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

	results[0].URL = "modifié"
	if batch.Results[0].URL != "https://exemple.fr" {
		t.Fatal("le batch doit avoir sa propre copie des résultats")
	}
}
