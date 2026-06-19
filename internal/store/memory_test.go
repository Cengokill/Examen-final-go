package store

import (
	"context"
	"errors"
	"testing"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// mockStore minimal pour tester l'interface Store sans réseau ni BDD.
type mockStore struct {
	data map[string]domain.Batch
}

func newMockStore() *mockStore {
	return &mockStore{data: make(map[string]domain.Batch)}
}

func (m *mockStore) Save(ctx context.Context, b domain.Batch) error {
	_ = ctx
	m.data[b.ID] = b
	return nil
}

func (m *mockStore) Get(ctx context.Context, id string) (domain.Batch, error) {
	_ = ctx
	b, ok := m.data[id]
	if !ok {
		return domain.Batch{}, domain.ErrBatchNotFound
	}
	return b, nil
}

func TestMemoryStoreTable(t *testing.T) {
	// test store ok
	cas := []struct {
		name      string
		setup     func(domain.Store)
		getID     string
		wantErr   error
		wantTotal int
	}{
		{
			name:    "lot introuvable",
			getID:   "b_absent",
			wantErr: domain.ErrBatchNotFound,
		},
		{
			name: "lot sauvegardé puis relu",
			setup: func(s domain.Store) {
				_ = s.Save(context.Background(), domain.NewBatch("b_ok", []domain.CheckResult{
					{URL: "https://a.test", Available: true, LatencyMs: 10},
				}))
			},
			getID:     "b_ok",
			wantTotal: 1,
		},
	}

	for _, tc := range cas {
		t.Run(tc.name, func(t *testing.T) {
			s := NewMemoryStore()
			if tc.setup != nil {
				tc.setup(s)
			}

			batch, err := s.Get(context.Background(), tc.getID)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					// METTRE errors.Is (sentinelle du cours)
					t.Fatalf("erreur attendue %v, obtenu %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if batch.Summary.Total != tc.wantTotal {
				t.Fatalf("total attendu %d, obtenu %d", tc.wantTotal, batch.Summary.Total)
			}
		})
	}
	// fmt.Println("MemoryStore table : Save/Get via interface Store")
}
