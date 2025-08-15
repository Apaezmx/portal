package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq" // Postgres driver
	"google.golang.org/grpc"

	pb "portal.com/portal/pkg/rag/v1" // The generated protobuf code
)

// config holds all the configuration for the service.
type config struct {
	grpcPort string
	dbConn   string
}

// server is used to implement rag.v1.RAGServiceServer.
type server struct {
	pb.UnimplementedRAGServiceServer
	db *sql.DB
}

// IndexContent implements rag.v1.RAGServiceServer
func (s *server) IndexContent(ctx context.Context, in *pb.IndexContentRequest) (*pb.IndexContentResponse, error) {
	log.Printf("Received IndexContent for URL: %v", in.Url)
	// TODO: Implement text chunking.
	// TODO: Implement vectorization (mock).
	// TODO: Store chunks and vectors in the database.
	log.Printf("Content length: %d", len(in.Content))
	return &pb.IndexContentResponse{}, nil
}

// RetrieveContext implements rag.v1.RAGServiceServer
func (s *server) RetrieveContext(ctx context.Context, in *pb.RetrieveContextRequest) (*pb.RetrieveContextResponse, error) {
	log.Printf("Received RetrieveContext for URL: %v with query: %s", in.Url, in.Query)
	// TODO: Implement similarity search in the database.
	// For now, return mock data.
	mockChunks := []string{
		"This is a mock chunk for the query.",
		"This is another relevant mock chunk.",
	}
	return &pb.RetrieveContextResponse{ContextChunks: mockChunks}, nil
}

func main() {
	var cfg config
	flag.StringVar(&cfg.grpcPort, "grpc-port", "50051", "The gRPC port to listen on")
	flag.StringVar(&cfg.dbConn, "db-conn", "host=postgres user=postgres password=postgres dbname=portal sslmode=disable", "PostgreSQL connection string")
	flag.Parse()

	// --- Database Connection ---
	db, err := sql.Open("postgres", cfg.dbConn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	// --- gRPC Server Setup ---
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterRAGServiceServer(s, &server{db: db})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
