package service_test

import (
	"context"
	"testing"
	"zdzira/backend/service"
	"zdzira/backend/store"
)

func newTestServices(t *testing.T) *service.Services {
	t.Helper()
	db, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}
	return service.New(store.New(db))
}

var ctx = context.Background()
