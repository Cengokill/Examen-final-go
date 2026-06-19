package store

import (
	"context"
	"sync"

	"github.com/Cengokill/Examen-final-go/internal/domain"
)

// MemoryStore persiste les lots en mémoire (map protégée par mutex).
type MemoryStore struct {
	mu      sync.RWMutex
	batches map[string]domain.Batch
}

// NewMemoryStore crée un store vide.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		batches: make(map[string]domain.Batch),
	}
}

// Save enregistre un lot.
func (s *MemoryStore) Save(ctx context.Context, b domain.Batch) error {
	_ = ctx
	s.mu.Lock()
	defer s.mu.Unlock()
	// fmt.Println("store Save : ", b.ID, " : ", len(b.Results), "résultats")
	s.batches[b.ID] = b
	return nil
}

// Get relit un lot par son identifiant.
func (s *MemoryStore) Get(ctx context.Context, id string) (domain.Batch, error) {
	_ = ctx
	s.mu.RLock()
	defer s.mu.RUnlock()

	batch, ok := s.batches[id]
	if !ok {
		// fmt.Println("store Get : lot introuvable", id)
		return domain.Batch{}, domain.ErrBatchNotFound
	}
	// fmt.Println("store Get OK : ", id, "total", batch.Summary.Total)
	return batch, nil
}
