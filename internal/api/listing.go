package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// FetchListing streams the listing.jsonl file and applies filters.
// It stops early once limit results are collected.
func (c *Client) FetchListing(sim, client string, limit int) ([]ListingEntry, error) {
	url := fmt.Sprintf("%s/%s/listing.jsonl", c.BaseURL, c.Suite)
	data, err := c.fetch(url, volatileTTL)
	if err != nil {
		return nil, err
	}

	var results []ListingEntry
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	// Increase scanner buffer for large lines
	scanner.Buffer(make([]byte, 0, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var entry ListingEntry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			continue // skip malformed lines
		}
		if sim != "" && !strings.Contains(strings.ToLower(entry.Name), strings.ToLower(sim)) {
			continue
		}
		if client != "" && !containsClient(entry.Clients, client) {
			continue
		}
		results = append(results, entry)
		if limit > 0 && len(results) >= limit {
			break
		}
	}

	return results, scanner.Err()
}

// FetchAllListing returns all listing entries matching sim filter (no limit).
func (c *Client) FetchAllListing(sim string) ([]ListingEntry, error) {
	return c.FetchListing(sim, "", 0)
}

func containsClient(clients []string, target string) bool {
	target = strings.ToLower(target)
	for _, cl := range clients {
		if strings.Contains(strings.ToLower(cl), target) {
			return true
		}
	}
	return false
}

// SortByTime sorts entries newest-first. Uses simple insertion sort since
// listing.jsonl entries are typically already roughly ordered.
func SortByTime(entries []ListingEntry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Start.After(entries[j-1].Start); j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}

// FormatTime returns a human-friendly relative time string.
func FormatTime(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		m := int(d.Minutes())
		return fmt.Sprintf("%dm ago", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		return fmt.Sprintf("%dh ago", h)
	default:
		days := int(d.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	}
}
