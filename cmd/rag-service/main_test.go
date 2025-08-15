package main

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	pb "portal.com/portal/pkg/rag/v1"
)

func TestIndexContent(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := &server{db: db}

	req := &pb.IndexContentRequest{
		Url:     "https://example.com",
		Content: "This is some test content.",
	}

	_, err = s.IndexContent(context.Background(), req)
	if err != nil {
		t.Errorf("IndexContent() error = %v, wantErr %v", err, false)
	}

	// This test is basic for now. When database logic is added to the IndexContent
	// method, we will add mock.ExpectExec(...) here to verify the SQL queries.
}

func TestRetrieveContext(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	s := &server{db: db}

	req := &pb.RetrieveContextRequest{
		Url:   "https://example.com",
		Query: "what is this?",
	}

	res, err := s.RetrieveContext(context.Background(), req)
	if err != nil {
		t.Errorf("RetrieveContext() error = %v, wantErr %v", err, false)
	}

	if len(res.ContextChunks) != 2 {
		t.Errorf("expected 2 mock chunks, got %d", len(res.ContextChunks))
	}

	// This test checks the current mock implementation. When database logic is added,
	// we will use mock.ExpectQuery(...) and mock.NewRows(...) to test the real logic.
}
