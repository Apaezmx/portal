package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"runtime"
	"time"

	"github.com/nats-io/nats.go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	expertpb "portal.com/portal/pkg/expert/v1"
)

// CrawledContentMessage is the structure of messages received from the crawler.
type CrawledContentMessage struct {
	URL       string    `json:"url"`
	Content   string    `json:"content"`
	CrawledAt time.Time `json:"crawled_at"`
}

// config holds all the configuration for the service.
type config struct {
	natsURL        string
	expertSvcAddr  string
	ragThreshold   int
}

func main() {
	var cfg config
	flag.StringVar(&cfg.natsURL, "nats-url", "nats://nats:4222", "The URL of the NATS server")
	flag.StringVar(&cfg.expertSvcAddr, "expert-svc-addr", "expert-service:50052", "The address of the Expert service")
	flag.IntVar(&cfg.ragThreshold, "rag-threshold", 4096, "The content length threshold to create a RAG expert")
	flag.Parse()

	// --- gRPC Client for Expert Service ---
	conn, err := grpc.NewClient(cfg.expertSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect to Expert service: %v", err)
	}
	defer conn.Close()
	expertSvcClient := expertpb.NewExpertServiceClient(conn)
	log.Println("Successfully connected to Expert service")

	// --- NATS Connection ---
	nc, err := nats.Connect(cfg.natsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	log.Println("Successfully connected to NATS")

	// --- NATS Subscription ---
	subject := "crawled-content"
	_, err = nc.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received a message on subject %s", subject)
		var contentMsg CrawledContentMessage
		if err := json.Unmarshal(msg.Data, &contentMsg); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			return // Don't process malformed messages
		}

		// Determine the expert type based on content length.
		expertType := expertpb.ExpertType_EXPERT_TYPE_SIMPLE
		if len(contentMsg.Content) > cfg.ragThreshold {
			expertType = expertpb.ExpertType_EXPERT_TYPE_RAG
		}

		// Call the Expert Service to process the content.
		req := &expertpb.CreateOrUpdateExpertRequest{
			Url:        contentMsg.URL,
			Content:    contentMsg.Content,
			ExpertType: expertType,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err := expertSvcClient.CreateOrUpdateExpert(ctx, req)
		if err != nil {
			log.Printf("Failed to call CreateOrUpdateExpert for URL %s: %v", contentMsg.URL, err)
			// In a real app, you might want to implement a retry mechanism or a dead-letter queue.
			return
		}

		log.Printf("Successfully processed and indexed URL: %s", contentMsg.URL)
	})
	if err != nil {
		log.Fatalf("failed to subscribe to subject %s: %v", subject, err)
	}

	log.Printf("Subscribed to subject '%s'", subject)

	// Keep the worker alive so the subscription can process messages.
	runtime.Goexit()
}
