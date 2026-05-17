package handler

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/mrckurz/CI-CD-MCM/internal/store"
)

func init() {
	sql.Register("fakehandler", &fakeHandlerDriver{})
}

type fakeHandlerDriver struct{}

type fakeHandlerConn struct {
	dsn string
}

type fakeHandlerStmt struct {
	query string
	dsn   string
}

type fakeHandlerRows struct {
	cols []string
	rows [][]driver.Value
	idx  int
}

type fakeHandlerResult struct {
	rowsAffected int64
}

func (*fakeHandlerDriver) Open(name string) (driver.Conn, error) {
	return &fakeHandlerConn{dsn: name}, nil
}

func (c *fakeHandlerConn) Prepare(query string) (driver.Stmt, error) {
	return &fakeHandlerStmt{query: query, dsn: c.dsn}, nil
}

func (*fakeHandlerConn) Close() error {
	return nil
}

func (*fakeHandlerConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported")
}

func (c *fakeHandlerConn) Ping(ctx context.Context) error {
	if strings.Contains(c.dsn, "fail=ping") {
		return errors.New("ping failed")
	}
	return nil
}

func (s *fakeHandlerStmt) Close() error {
	return nil
}

func (s *fakeHandlerStmt) NumInput() int {
	return -1
}

func (s *fakeHandlerStmt) Exec(args []driver.Value) (driver.Result, error) {
	q := strings.TrimSpace(s.query)
	if strings.Contains(s.dsn, "fail=update") {
		return nil, errors.New("update failed")
	}
	if strings.Contains(s.dsn, "fail=delete") {
		return nil, errors.New("delete failed")
	}
	if strings.HasPrefix(q, "UPDATE products SET name = $1, price = $2 WHERE id = $3") {
		id := toInt(args[2])
		if id == 999 {
			return fakeHandlerResult{rowsAffected: 0}, nil
		}
		return fakeHandlerResult{rowsAffected: 1}, nil
	}
	if strings.HasPrefix(q, "DELETE FROM products WHERE id = $1") {
		id := toInt(args[0])
		if id == 999 {
			return fakeHandlerResult{rowsAffected: 0}, nil
		}
		return fakeHandlerResult{rowsAffected: 1}, nil
	}
	return nil, errors.New("unexpected exec")
}

func (s *fakeHandlerStmt) Query(args []driver.Value) (driver.Rows, error) {
	q := strings.TrimSpace(s.query)
	if strings.Contains(s.dsn, "fail=getall") {
		return nil, errors.New("select failed")
	}
	if strings.HasPrefix(q, "SELECT id, name, price FROM products ORDER BY id") {
		return &fakeHandlerRows{cols: []string{"id", "name", "price"}, rows: [][]driver.Value{{int64(1), "brodukt", 2.00}}}, nil
	}
	if strings.HasPrefix(q, "SELECT id, name, price FROM products WHERE id = $1") {
		id := toInt(args[0])
		if id == 999 {
			return &fakeHandlerRows{cols: []string{"id", "name", "price"}, rows: [][]driver.Value{}}, nil
		}
		return &fakeHandlerRows{cols: []string{"id", "name", "price"}, rows: [][]driver.Value{{int64(1), "brodukt", 2.00}}}, nil
	}
	if strings.HasPrefix(q, "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id") {
		return &fakeHandlerRows{cols: []string{"id"}, rows: [][]driver.Value{{int64(1)}}}, nil
	}
	return nil, errors.New("unexpected query")
}

func (r *fakeHandlerRows) Columns() []string {
	return r.cols
}

func (r *fakeHandlerRows) Close() error {
	return nil
}

func (r *fakeHandlerRows) Next(dest []driver.Value) error {
	if r.idx >= len(r.rows) {
		return io.EOF
	}
	for i := range dest {
		dest[i] = r.rows[r.idx][i]
	}
	r.idx++
	return nil
}

func (r fakeHandlerResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (r fakeHandlerResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func toInt(v driver.Value) int {
	switch x := v.(type) {
	case int64:
		return int(x)
	case int:
		return x
	}
	return 0
}

func setupPostgresRouter(t *testing.T, dsn string) (*mux.Router, *PostgresHandler) {
	db, err := sql.Open("fakehandler", dsn)
	if err != nil {
		t.Fatalf("failed to open fake handler db: %v", err)
	}
	s := &store.PostgresStore{DB: db}
	h := NewPostgresHandler(s)
	r := mux.NewRouter()
	h.RegisterRoutes(r)
	return r, h
}

func TestPostgresHandlerHealthSuccess(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresHandlerGetProducts(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresHandlerGetProductNotFound(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("GET", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresHandlerGetProductSuccess(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("GET", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresHandlerCreateProduct(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Neichs Produkt","price":2.01}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rr.Code)
	}
}

func TestPostgresHandlerUpdateProductSuccess(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{"name":"Updates Produkt","price":2.00}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresHandlerUpdateProductNotFound(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("PUT", "/products/999", strings.NewReader(`{"name":"Ned gfunden","price":4.04}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresHandlerDeleteProductSuccess(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("DELETE", "/products/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostgresHandlerDeleteProductNotFound(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("DELETE", "/products/999", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestPostgresHandlerHealthFailure(t *testing.T) {
	r, _ := setupPostgresRouter(t, "fail=ping")

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rr.Code)
	}
}

func TestPostgresHandlerGetProductsError(t *testing.T) {
	r, _ := setupPostgresRouter(t, "fail=getall")

	req := httptest.NewRequest("GET", "/products", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rr.Code)
	}
}

func TestPostgresHandlerCreateProductInvalidJSON(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"Widget","price":}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresHandlerCreateProductInvalidProduct(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("POST", "/products", strings.NewReader(`{"name":"","price":4.00}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestPostgresHandlerUpdateProductInvalidJSON(t *testing.T) {
	r, _ := setupPostgresRouter(t, "")

	req := httptest.NewRequest("PUT", "/products/1", strings.NewReader(`{"name":"Ka göd","price":}`))
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
