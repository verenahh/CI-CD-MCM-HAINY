package store

import (
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

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

func TestCreateAndGet(t *testing.T) {
	s := NewMemoryStore()

	created := s.Create(model.Product{Name: "Neichs produkt", Price: 9.99})
	if created.ID != 1 {
		t.Fatalf("expected created product ID 1, got %d", created.ID)
	}

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "Neichs produkt" || got.Price != 9.99 {
		t.Fatalf("unexpected product retrieved: %#v", got)
	}
}

func TestUpdateExistingProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Neichs produkt", Price: 9.99})

	updated, err := s.Update(created.ID, model.Product{Name: "Preis auffe", Price: 19.99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if updated.ID != created.ID || updated.Name != "Preis auffe" || updated.Price != 19.99 {
		t.Fatalf("unexpected updated product: %#v", updated)
	}
}

func TestUpdateNonExistentProduct(t *testing.T) {
	s := NewMemoryStore()
	_, err := s.Update(999, model.Product{Name: "Preis auffe", Price: 19.99})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteExistingProduct(t *testing.T) {
	s := NewMemoryStore()
	created := s.Create(model.Product{Name: "Neichs produkt", Price: 9.99})

	err := s.Delete(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = s.GetByID(created.ID)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
