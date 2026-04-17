package utils

import (
	"context"
	"errors"
	"testing"
	"time"
)

// MockStore para pruebas
type MockStore struct {
	getFunc func(ctx context.Context, key string, dest interface{}) error
	setFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	delFunc func(ctx context.Context, key string) error
	preFunc func(ctx context.Context, prefix string) error
}

func (m *MockStore) Get(ctx context.Context, key string, dest interface{}) error {
	return m.getFunc(ctx, key, dest)
}
func (m *MockStore) Set(ctx context.Context, key string, value interface{}, exp time.Duration) error {
	return m.setFunc(ctx, key, value, exp)
}
func (m *MockStore) Delete(ctx context.Context, key string) error { return m.delFunc(ctx, key) }
func (m *MockStore) DeleteByPrefix(ctx context.Context, prefix string) error {
	return m.preFunc(ctx, prefix)
}

func TestGetOrSet(t *testing.T) {
	ctx := context.Background()

	t.Run("Cache Miss y Fetch exitoso", func(t *testing.T) {
		mock := &MockStore{
			getFunc: func(ctx context.Context, key string, dest interface{}) error {
				return errors.New("not found")
			},
			setFunc: func(ctx context.Context, key string, value interface{}, exp time.Duration) error {
				return nil
			},
		}
		cm := NewCacheManager(mock)

		called := false
		fetch := func() (string, error) {
			called = true
			return "data", nil
		}

		res, err := GetOrSet(ctx, cm, "key", time.Hour, fetch)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if res != "data" {
			t.Errorf("esperado 'data', obtenido '%s'", res)
		}
		if !called {
			t.Error("fetch no fue llamado")
		}
	})

	t.Run("Cache Hit", func(t *testing.T) {
		mock := &MockStore{
			getFunc: func(ctx context.Context, key string, dest interface{}) error {
				*dest.(*string) = "cached"
				return nil
			},
		}
		cm := NewCacheManager(mock)

		fetch := func() (string, error) {
			t.Error("fetch no debería ser llamado")
			return "data", nil
		}

		res, err := GetOrSet(ctx, cm, "key", time.Hour, fetch)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if res != "cached" {
			t.Errorf("esperado 'cached', obtenido '%s'", res)
		}
	})

	t.Run("Store nil (No cache mode)", func(t *testing.T) {
		cm := NewCacheManager(nil)

		called := false
		fetch := func() (string, error) {
			called = true
			return "live", nil
		}

		res, err := GetOrSet(ctx, cm, "key", time.Hour, fetch)
		if err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if res != "live" {
			t.Errorf("esperado 'live', obtenido '%s'", res)
		}
		if !called {
			t.Error("fetch no fue llamado")
		}
	})
}
