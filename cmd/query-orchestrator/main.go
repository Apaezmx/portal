package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	expertpb "portal.com/portal/pkg/expert/v1"
	pb "portal.com/portal/pkg/orchestrator/v1"
)

// config holds all the configuration for the service.
type config struct {
	grpcPort      string
	expertSvcAddr string
}

// server implements the QueryOrchestratorService.
type server struct {
	pb.UnimplementedQueryOrchestratorServiceServer
	expertSvcClient expertpb.ExpertServiceClient
}

// Search implements orchestrator.v1.QueryOrchestratorServiceServer
func (s *server) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	log.Printf("Received Search request with query: %s", in.Query)

	// TODO: This is simplified routing. A real implementation would have a
	// dynamic way to select which experts to query based on the user's query.
	// For now, we'll hardcode a single expert to consult.
	expertURL := "http://gocolly.dev/" // The same URL our crawler starts with

	log.Printf("Querying expert for URL: %s", expertURL)

	expertReq := &expertpb.QueryExpertRequest{
		Url:   expertURL,
		Query: in.Query,
	}

	expertRes, err := s.expertSvcClient.QueryExpert(ctx, expertReq)
	if err != nil {
		log.Printf("Failed to query expert service for URL %s: %v", expertURL, err)
		return nil, fmt.Errorf("expert query failed: %w", err)
	}

	// TODO: This is simplified synthesis. A real implementation would use an LLM
	// to synthesize answers from multiple experts into a coherent summary.
	summary := expertRes.Answer

	source := &pb.Source{
		Url:     expertURL,
		Title:   "Mock Title", // TODO: Get title from expert/metadata
		Snippet: expertRes.Answer,
	}

	return &pb.SearchResponse{
		Summary: summary,
		Sources: []*pb.Source{source},
	}, nil
}

func main() {
	var cfg config
	flag.StringVar(&cfg.grpcPort, "grpc-port", "50053", "The gRPC port to listen on")
	flag.StringVar(&cfg.expertSvcAddr, "expert-svc-addr", "expert-service:50052", "The address of the Expert service")
	flag.Parse()

	// --- gRPC Client for Expert Service ---
	conn, err := grpc.NewClient(cfg.expertSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to Expert service: %v", err)
	}
	defer conn.Close()
	expertSvcClient := expertpb.NewExpertServiceClient(conn)
	log.Println("Successfully connected to Expert service")

	// --- gRPC Server Setup for this service ---
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterQueryOrchestratorServiceServer(s, &server{expertSvcClient: expertSvcClient})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
