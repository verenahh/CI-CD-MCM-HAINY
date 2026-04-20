package store

import (
	"errors"
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()
	p := model.Product{Name: "Apple", Price: 1599}

	created := s.Create(p)

	found, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found.Name != p.Name || found.Price != p.Price {
		t.Errorf("expected product %v, got %v", p, found)
	}
}

func TestGetAllEmpty(t *testing.T) {
	s := NewMemoryStore()
	products := s.GetAll()
	if len(products) != 0 {
		t.Errorf("expected 0 products, got %d", len(products))
	}
}

func TestDeleteNonExistent(t *testing.T) {
	s := NewMemoryStore()
	err := s.Delete(999)
	if err != ErrNotFound {
		t.Error("expected ErrNotFound when deleting non-existent product")
	}
}

func TestUpdateProduct(t *testing.T) {
	s := NewMemoryStore()
	p := model.Product{Name: "Pear", Price: 18.99}

	created := s.Create(p)

	updated := model.Product{Name: "Peach", Price: 100}
	result, err := s.Update(created.ID, updated)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	if result.Name != "Peach" || result.Price != 100 {
		t.Errorf("Update returned incorrect values: got %v, want name=Peach, price=100", result)
	}

	retrieved, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("GetByID failed after update: %v", err)
	}
	if retrieved.Name != "Peach" || retrieved.Price != 100 {
		t.Errorf("Product not updated in store: got %v, want name=Peach, price=100", retrieved)
	}
}

func TestDeleteProduct(t *testing.T) {
	s := NewMemoryStore()
	p := model.Product{Name: "Raspberry", Price: 53.99}

	created := s.Create(p)

	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = s.GetByID(created.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected ErrNotFound after deletion, got %v", err)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		setup     func(*MemoryStore) // optional setup function
		wantError error
	}{
		{
			name:      "non-existent ID in empty store",
			id:        1,
			wantError: ErrNotFound,
		},
		{
			name:      "negative ID",
			id:        -1,
			wantError: ErrNotFound,
		},
		{
			name:      "zero ID",
			id:        0,
			wantError: ErrNotFound,
		},
		{
			name: "non-existent ID when store has products",
			id:   9999,
			setup: func(s *MemoryStore) {
				s.Create(model.Product{Name: "Paspberry Rie", Price: 1899})
			},
			wantError: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMemoryStore()
			if tt.setup != nil {
				tt.setup(s)
			}

			_, err := s.GetByID(tt.id)
			if !errors.Is(err, tt.wantError) {
				t.Errorf("GetByID(%d) returned error %v, want %v", tt.id, err, tt.wantError)
			}
		})
	}
}

//All puns intended
