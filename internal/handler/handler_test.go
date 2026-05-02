package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func setupRouter() (*mux.Router, *Handler) {
	s := store.NewMemoryStore()
	h := NewHandler(s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r, h
}

func TestHealthEndpoint(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductsEmpty(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestCreateAndGetProduct(t *testing.T) {
	r, _ := setupRouter()

	// Create
	body := `{"name":"Widget","price":9.99}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", rr.Code)
	}

	// Get
	req = httptest.NewRequest("GET", "/products/1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}
}

func TestGetProductNotFound(t *testing.T) {
	r, _ := setupRouter()

	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rr.Code)
	}
}

// TODO: Add tests for UpdateProduct, DeleteProduct, and invalid payloads

func TestUpdateProduct(t *testing.T) {
	r, _ := setupRouter()

	createBody := `{"name":"Outdated Product","price":55.55}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(createBody))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	updateBody := `{"name":"Updated Product","price":100.00}`
	req = httptest.NewRequest("PUT", "/products/1", strings.NewReader(updateBody))
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "Updated Product") {
		t.Errorf("expected body to contain 'Updated Product', got %s", rr.Body.String())
	}
}

func TestDeleteProduct(t *testing.T) {
	r, _ := setupRouter()

	createBody := `{"name":"Deleteable Product","price":55.55}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(createBody))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	req = httptest.NewRequest("DELETE", "/products/1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	req = httptest.NewRequest("GET", "/products/1", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 for deleted product, got %d", rr.Code)
	}
}

func TestCreateInvalidProduct(t *testing.T) {
	r, _ := setupRouter()

	invalidBody := `{"name":"","price":4.00}`
	req := httptest.NewRequest("POST", "/products", strings.NewReader(invalidBody))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid payload, got %d", rr.Code)
	}
}
