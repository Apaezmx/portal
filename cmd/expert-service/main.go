package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "portal.com/portal/pkg/expert/v1" // This service's generated code
	ragpb "portal.com/portal/pkg/rag/v1" // RAG service's generated code
)

// config holds all the configuration for the service.
type config struct {
	grpcPort   string
	dbConn     string
	ragSvcAddr string
}

// server implements the ExpertService.
type server struct {
	pb.UnimplementedExpertServiceServer
	db           *sql.DB
	ragSvcClient ragpb.RAGServiceClient
}

// CreateOrUpdateExpert implements expert.v1.ExpertServiceServer
func (s *server) CreateOrUpdateExpert(ctx context.Context, in *pb.CreateOrUpdateExpertRequest) (*pb.CreateOrUpdateExpertResponse, error) {
	log.Printf("Received CreateOrUpdateExpert for URL: %v", in.Url)
	// TODO: Check if expert exists in DB.
	// TODO: If it's a RAG expert, call the RAG service's IndexContent method.
	// s.ragSvcClient.IndexContent(...)
	// TODO: Insert/update expert metadata in the 'experts' table.
	return &pb.CreateOrUpdateExpertResponse{ExpertId: "mock-expert-id"}, nil
}

// QueryExpert implements expert.v1.ExpertServiceServer
func (s *server) QueryExpert(ctx context.Context, in *pb.QueryExpertRequest) (*pb.QueryExpertResponse, error) {
	log.Printf("Received QueryExpert for URL: %v", in.Url)
	// TODO: Fetch expert metadata from DB to see if it's simple or RAG.
	// TODO: If RAG, call RAG service's RetrieveContext method.
	// TODO: If simple, get content from DB.
	// TODO: Call LLM (mock for now) with context and query.
	return &pb.QueryExpertResponse{Answer: "This is a mock answer from the expert."}, nil
}

func main() {
	var cfg config
	flag.StringVar(&cfg.grpcPort, "grpc-port", "50052", "The gRPC port to listen on")
	flag.StringVar(&cfg.dbConn, "db-conn", "host=postgres user=postgres password=postgres dbname=portal sslmode=disable", "PostgreSQL connection string")
	flag.StringVar(&cfg.ragSvcAddr, "rag-svc-addr", "rag-service:50051", "The address of the RAG service")
	flag.Parse()

	// --- Database Connection ---
	db, err := sql.Open("postgres", cfg.dbConn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to the database")

	// --- gRPC Client for RAG Service ---
	conn, err := grpc.NewClient(cfg.ragSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to RAG service: %v", err)
	}
	defer conn.Close()
	ragSvcClient := ragpb.NewRAGServiceClient(conn)
	log.Println("Successfully connected to RAG service")

	// --- gRPC Server Setup for this service ---
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterExpertServiceServer(s, &server{db: db, ragSvcClient: ragSvcClient})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
