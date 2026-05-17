package store

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/mrckurz/CI-CD-MCM/internal/model"
)

func init() {
	sql.Register("fake", &fakeDriver{})
}

type fakeDriver struct{}

type fakeConn struct{}

type fakeStmt struct {
	query string
}

type fakeRows struct {
	cols []string
	rows [][]driver.Value
	idx  int
}

type fakeResult struct {
	lastInsertID int64
	rowsAffected int64
}

func (*fakeDriver) Open(name string) (driver.Conn, error) {
	return &fakeConn{}, nil
}

func (*fakeConn) Prepare(query string) (driver.Stmt, error) {
	return &fakeStmt{query: query}, nil
}

func (*fakeConn) Close() error {
	return nil
}

func (*fakeConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions not supported")
}

func (s *fakeStmt) Close() error {
	return nil
}

func (s *fakeStmt) NumInput() int {
	return -1
}

func (s *fakeStmt) Exec(args []driver.Value) (driver.Result, error) {
	q := strings.TrimSpace(s.query)
	if strings.HasPrefix(q, "CREATE TABLE IF NOT EXISTS products") {
		return fakeResult{lastInsertID: 0, rowsAffected: 0}, nil
	}
	if strings.HasPrefix(q, "UPDATE products SET name = $1, price = $2 WHERE id = $3") {
		id := toInt(args[2])
		if id == 999 {
			return fakeResult{rowsAffected: 0}, nil
		}
		return fakeResult{rowsAffected: 1}, nil
	}
	if strings.HasPrefix(q, "DELETE FROM products WHERE id = $1") {
		id := toInt(args[0])
		if id == 999 {
			return fakeResult{rowsAffected: 0}, nil
		}
		return fakeResult{rowsAffected: 1}, nil
	}
	return nil, errors.New("unexpected exec")
}

func (s *fakeStmt) Query(args []driver.Value) (driver.Rows, error) {
	q := strings.TrimSpace(s.query)
	if strings.HasPrefix(q, "SELECT id, name, price FROM products ORDER BY id") {
		return &fakeRows{cols: []string{"id", "name", "price"}, rows: [][]driver.Value{{int64(1), "Widget", 9.99}}}, nil
	}
	if strings.HasPrefix(q, "SELECT id, name, price FROM products WHERE id = $1") {
		id := toInt(args[0])
		if id == 999 {
			return &fakeRows{cols: []string{"id", "name", "price"}, rows: [][]driver.Value{}}, nil
		}
		return &fakeRows{cols: []string{"id", "name", "price"}, rows: [][]driver.Value{{int64(1), "Widget", 9.99}}}, nil
	}
	if strings.HasPrefix(q, "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id") {
		return &fakeRows{cols: []string{"id"}, rows: [][]driver.Value{{int64(1)}}}, nil
	}
	return nil, errors.New("unexpected query")
}

func (r *fakeRows) Columns() []string {
	return r.cols
}

func (r *fakeRows) Close() error {
	return nil
}

func (r *fakeRows) Next(dest []driver.Value) error {
	if r.idx >= len(r.rows) {
		return io.EOF
	}
	for i := range dest {
		dest[i] = r.rows[r.idx][i]
	}
	r.idx++
	return nil
}

func (r fakeResult) LastInsertId() (int64, error) {
	return r.lastInsertID, nil
}

func (r fakeResult) RowsAffected() (int64, error) {
	return r.rowsAffected, nil
}

func toInt(v driver.Value) int {
	switch x := v.(type) {
	case int64:
		return int(x)
	case int:
		return x
	case string:
		return int(x[0])
	}
	return 0
}

func TestNewPostgresStoreInvalidPort(t *testing.T) {
	_, err := NewPostgresStore("127.0.0.1", "1", "catalog", "catalog123", "productcatalog")
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestPostgresEnsureTable(t *testing.T) {
	db, err := sql.Open("fake", "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	s := &PostgresStore{DB: db}
	if err := s.EnsureTable(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPostgresGetAll(t *testing.T) {
	db, err := sql.Open("fake", "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	s := &PostgresStore{DB: db}
	products, err := s.GetAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 1 || products[0].ID != 1 {
		t.Fatalf("unexpected products: %#v", products)
	}
}

func TestPostgresGetByIDNotFound(t *testing.T) {
	db, err := sql.Open("fake", "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	s := &PostgresStore{DB: db}
	_, err = s.GetByID(999)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresCreate(t *testing.T) {
	db, err := sql.Open("fake", "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	s := &PostgresStore{DB: db}
	p, err := s.Create(model.Product{Name: "Widget", Price: 9.99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 1 {
		t.Fatalf("expected ID 1, got %d", p.ID)
	}
}

func TestPostgresUpdateNotFound(t *testing.T) {
	db, err := sql.Open("fake", "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	s := &PostgresStore{DB: db}
	_, err = s.Update(999, model.Product{Name: "Widget", Price: 9.99})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresUpdateFound(t *testing.T) {
	db, err := sql.Open("fake", "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	s := &PostgresStore{DB: db}
	p, err := s.Update(1, model.Product{Name: "Updated", Price: 19.99})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.ID != 1 {
		t.Fatalf("expected ID 1, got %d", p.ID)
	}
}

func TestPostgresDeleteNotFound(t *testing.T) {
	db, err := sql.Open("fake", "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	s := &PostgresStore{DB: db}
	err = s.Delete(999)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresDeleteFound(t *testing.T) {
	db, err := sql.Open("fake", "")
	if err != nil {
		t.Fatalf("failed to open fake db: %v", err)
	}
	s := &PostgresStore{DB: db}
	if err := s.Delete(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
