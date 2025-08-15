package main

import (
	"encoding/json"
	"flag"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/nats-io/nats.go"
)

// CrawledContentMessage defines the structure of the message sent to NATS.
type CrawledContentMessage struct {
	URL       string    `json:"url"`
	Content   string    `json:"content"`
	CrawledAt time.Time `json:"crawled_at"`
}

// config holds all the configuration for the service.
type config struct {
	natsURL        string
	allowedDomains string
	startURL       string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.natsURL, "nats-url", "nats://nats:4222", "The URL of the NATS server")
	flag.StringVar(&cfg.allowedDomains, "allowed-domains", "gocolly.dev", "A comma-separated list of domains to allow crawling")
	flag.StringVar(&cfg.startURL, "start-url", "http://gocolly.dev/", "The initial URL to start crawling from")
	flag.Parse()

	// --- NATS Connection ---
	nc, err := nats.Connect(cfg.natsURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()
	log.Println("Successfully connected to NATS")

	// Instantiate default collector
	c := colly.NewCollector(
		colly.AllowedDomains(strings.Split(cfg.allowedDomains, ",")...),
	)

	// publisher creates and sends a message to the NATS queue.
	publisher := func(url, content string) {
		if content == "" {
			log.Printf("Skipping empty content for URL: %s", url)
			return
		}

		msg := CrawledContentMessage{
			URL:       url,
			Content:   content,
			CrawledAt: time.Now().UTC(),
		}

		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Printf("failed to marshal message for url %s: %v", url, err)
			return
		}

		// Publish the message to the "crawled-content" subject.
		if err := nc.Publish("crawled-content", msgBytes); err != nil {
			log.Printf("failed to publish message for url %s: %v", url, err)
		} else {
			log.Printf("Published content for URL: %s", url)
		}
	}

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	// When a page is scraped, extract its text and publish it.
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// This is a naive content extraction. A real implementation would use
		// a library like go-readability to get only the main article text.
		publisher(e.Request.URL.String(), e.Text)
	})

	log.Printf("Starting crawl at %s", cfg.startURL)
	if err := c.Visit(cfg.startURL); err != nil {
		log.Fatalf("failed to start crawl: %v", err)
	}

	// The collector runs asynchronously, so we block forever.
	// In a real app, you'd handle signals for graceful shutdown.
	select {}
}
