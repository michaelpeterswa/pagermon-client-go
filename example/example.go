package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	pagermon "github.com/michaelpeterswa/go-lib"
	"github.com/michaelpeterswa/go-lib/multimonng"
)

const (
	apiKeyKey     = "PAGERMON_API_KEY"
	baseURLKey    = "PAGERMON_BASE_URL"
	identifierKey = "PAGERMON_IDENTIFIER"

	exampleMultimonNGRawMessage = "POCSAG1200: Address: 1234567  Function: 0  Alpha:   Aid - Emergency; *FTAC - 1*;  Test Emergency Location; 7xxx Test Rd NE, RM; A1; 47.0;-122.0<EOT><NUL>"
)

func main() {

	ctx := context.Background()

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	apiKey := os.Getenv(apiKeyKey)
	if apiKey == "" {
		panic("PAGERMON_API_KEY environment variable must be set")
	}

	baseURL := os.Getenv(baseURLKey)
	if baseURL == "" {
		panic("PAGERMON_BASE_URL environment variable must be set")
	}

	identifier := os.Getenv(identifierKey)
	if identifier == "" {
		panic("PAGERMON_IDENTIFIER environment variable must be set")
	}

	// Create a new PagerMon client
	pagerMonClient := pagermon.NewPagerMonClient(httpClient, apiKey, baseURL)

	// Parse the raw multimon-ng message
	multimonNGMessage, err := multimonng.ParseMultimonLine(exampleMultimonNGRawMessage)
	if err != nil {
		panic(fmt.Errorf("error parsing multimon-ng message: %w", err))
	}

	err = pagerMonClient.SendMessage(ctx, pagermon.MultimonNGMessageToPagerMonMessage(multimonNGMessage, identifier))
	if err != nil {
		panic(fmt.Errorf("error sending message to pagermon: %w", err))
	}
}
