package websearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func prepareFirecrawlSearchRequest(ctx context.Context, apikey string, searchRequestBody FirecrawlSearchRequest) (*http.Request, error) {
	bodyBytes, err := json.Marshal(searchRequestBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling search request body: %w", err)

	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/search", firecrawlApiBaseUri), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func prepareFirecrawlScrapRequest(ctx context.Context, apikey string, scrapRequest FirecrawlScrapRequest) (*http.Request, error) {
	bodyBytes, err := json.Marshal(scrapRequest)
	if err != nil {
		return nil, fmt.Errorf("marshalling search request body: %w", err)

	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/scrape", firecrawlApiBaseUri), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apikey))
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func parseDomainQueries(answer string) (QueryDomains, error) {
	var queries QueryDomains
	err := json.Unmarshal([]byte(answer), &queries)
	if err != nil {
		return QueryDomains{}, fmt.Errorf("failed to parse domain queries: %w", err)
	}
	return queries, nil
}
