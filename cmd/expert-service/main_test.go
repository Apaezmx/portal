package main

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"google.golang.org/grpc"

	pb "portal.com/portal/pkg/expert/v1"
	ragpb "portal.com/portal/pkg/rag/v1"
)

// mockRAGServiceClient is a mock implementation of RAGServiceClient.
type mockRAGServiceClient struct {
	ragpb.UnimplementedRAGServiceClient
	// This mock can be extended to control responses for different test cases.
}

// IndexContent is the mock implementation for the RAG service's IndexContent method.
func (m *mockRAGServiceClient) IndexContent(ctx context.Context, in *ragpb.IndexContentRequest, opts ...grpc.CallOption) (*ragpb.IndexContentResponse, error) {
	// In a real test, you might check `in` and return different responses.
	return &ragpb.IndexContentResponse{}, nil
}

// RetrieveContext is the mock implementation for the RAG service's RetrieveContext method.
func (m *mockRAGServiceClient) RetrieveContext(ctx context.Context, in *ragpb.RetrieveContextRequest, opts ...grpc.CallOption) (*ragpb.RetrieveContextResponse, error) {
	return &ragpb.RetrieveContextResponse{ContextChunks: []string{"mocked chunk"}}, nil
}

func TestCreateOrUpdateExpert(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRagClient := &mockRAGServiceClient{}

	s := &server{
		db:           db,
		ragSvcClient: mockRagClient,
	}

	req := &pb.CreateOrUpdateExpertRequest{
		Url:        "https://example.com",
		Content:    "test content",
		ExpertType: pb.ExpertType_EXPERT_TYPE_SIMPLE,
	}

	res, err := s.CreateOrUpdateExpert(context.Background(), req)
	if err != nil {
		t.Errorf("CreateOrUpdateExpert() error = %v, wantErr %v", err, false)
	}

	if res.ExpertId != "mock-expert-id" {
		t.Errorf("expected mock expert id, got %s", res.ExpertId)
	}
}

func TestQueryExpert(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mockRagClient := &mockRAGServiceClient{}

	s := &server{
		db:           db,
		ragSvcClient: mockRagClient,
	}

	req := &pb.QueryExpertRequest{
		Url:   "https://example.com",
		Query: "what is this?",
	}

	res, err := s.QueryExpert(context.Background(), req)
	if err != nil {
		t.Errorf("QueryExpert() error = %v, wantErr %v", err, false)
	}

	if res.Answer != "This is a mock answer from the expert." {
		t.Errorf("unexpected mock answer: %s", res.Answer)
	}
}
