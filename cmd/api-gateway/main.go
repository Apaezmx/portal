package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	expertpb "portal.com/portal/pkg/expert/v1"
	orchpb "portal.com/portal/pkg/orchestrator/v1"
)

// config holds all the configuration for the service.
type config struct {
	httpPort      string
	expertSvcAddr string
	orchSvcAddr   string
}

// apiServer holds the clients for the backend gRPC services.
type apiServer struct {
	expertSvcClient expertpb.ExpertServiceClient
	orchSvcClient   orchpb.QueryOrchestratorServiceClient
}

// searchHandler handles requests to the /search endpoint.
func (s *apiServer) searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req orchpb.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grpcRes, err := s.orchSvcClient.Search(r.Context(), &req)
	if err != nil {
		// In a real app, inspect gRPC error code for better HTTP status mapping.
		http.Error(w, "backend service error", http.StatusInternalServerError)
		log.Printf("Error from orchestrator service: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grpcRes)
}

// expertHandler handles requests to the /e/{url} endpoint.
func (s *apiServer) expertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract URL from path, e.g., /e/https://example.com
	url := strings.TrimPrefix(r.URL.Path, "/e/")
	if url == "" {
		http.Error(w, "URL path parameter is missing", http.StatusBadRequest)
		return
	}

	var req expertpb.QueryExpertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Url = url // Set the URL from the path

	grpcRes, err := s.expertSvcClient.QueryExpert(r.Context(), &req)
	if err != nil {
		http.Error(w, "backend service error", http.StatusInternalServerError)
		log.Printf("Error from expert service: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(grpcRes)
}

func main() {
	var cfg config
	flag.StringVar(&cfg.httpPort, "http-port", "8080", "The HTTP port to listen on")
	flag.StringVar(&cfg.expertSvcAddr, "expert-svc-addr", "expert-service:50052", "The address of the Expert service")
	flag.StringVar(&cfg.orchSvcAddr, "orch-svc-addr", "query-orchestrator:50053", "The address of the Query Orchestrator service")
	flag.Parse()

	// --- gRPC Client for Expert Service ---
	expertConn, err := grpc.NewClient(cfg.expertSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to Expert service: %v", err)
	}
	defer expertConn.Close()
	expertSvcClient := expertpb.NewExpertServiceClient(expertConn)
	log.Println("Successfully connected to Expert service")

	// --- gRPC Client for Query Orchestrator Service ---
	orchConn, err := grpc.NewClient(cfg.orchSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to Query Orchestrator service: %v", err)
	}
	defer orchConn.Close()
	orchSvcClient := orchpb.NewQueryOrchestratorServiceClient(orchConn)
	log.Println("Successfully connected to Query Orchestrator service")

	// --- HTTP Server Setup ---
	server := &apiServer{
		expertSvcClient: expertSvcClient,
		orchSvcClient:   orchSvcClient,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/search", server.searchHandler)
	mux.HandleFunc("/e/", server.expertHandler)

	// Serve the frontend files
	fs := http.FileServer(http.Dir("./frontend"))
	mux.Handle("/", fs)

	log.Printf("API Gateway listening on :%s", cfg.httpPort)
	if err := http.ListenAndServe(":"+cfg.httpPort, mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
