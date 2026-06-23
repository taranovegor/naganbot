package drand

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// quicknet: 3-second rounds, unchained, BLS12-381
const chainHash = "52db9ba70e0cc0f6eaf7803dd07447a1f5477735fd3f661792ba94600c84e971"

var relayURLs = []string{
	"https://api.drand.sh",
	"https://api2.drand.sh",
	"https://api3.drand.sh",
	"https://drand.cloudflare.com",
}

type Beacon struct {
	Round      uint64 `json:"round"`
	Randomness string `json:"randomness"`
	Signature  string `json:"signature"`
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	chainHash  string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		chainHash:  chainHash,
	}
}

func NewClientWithURL(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		baseURL:    baseURL,
		chainHash:  chainHash,
	}
}

func (c *Client) GetLatest(ctx context.Context) (*Beacon, error) {
	urls := relayURLs
	if c.baseURL != "" {
		urls = []string{c.baseURL}
	}

	var lastErr error
	for _, base := range urls {
		url := fmt.Sprintf("%s/%s/public/latest", base, c.chainHash)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			lastErr = err
			continue
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			lastErr = fmt.Errorf("relay %s returned status %d", base, resp.StatusCode)
			continue
		}

		var beacon Beacon
		err = json.NewDecoder(resp.Body).Decode(&beacon)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		if beacon.Randomness == "" {
			lastErr = fmt.Errorf("relay %s returned empty randomness", base)
			continue
		}

		return &beacon, nil
	}

	return nil, fmt.Errorf("all drand relays failed: %w", lastErr)
}

func (c *Client) ProofURL(round uint64, gameID uuid.UUID) string {
	return fmt.Sprintf("%s/%s/public/%d?game[id]=%s", relayURLs[0], c.chainHash, round, gameID.String())
}
