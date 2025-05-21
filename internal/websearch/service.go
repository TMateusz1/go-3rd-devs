package websearch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TMateusz1/go-3rd-devs/internal/ai"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

const firecrawlApiBaseUri = "https://api.firecrawl.dev/v1"

type Service interface {
	IsSearchRequired(context.Context, string) bool
	GetDomainQueries(context.Context, string, []AllowedDomain) (QueryDomains, error)
	SearchForSpecificPages(context.Context, QueryDomains) ([]SearchResult, error)
	ScoreResults(context.Context, []SearchResult, string) ([]WebPage, error)
	ScrapWebpages(context.Context, []WebPage, string) ([]ScrappedWebPage, error)
}

const (
	doSearch string = "1"
)

func NewService(as ai.Service) (Service, error) {
	firecrawlApiKey, ok := os.LookupEnv("FIRECRAWL_API_KEY")
	if !ok {
		return nil, fmt.Errorf("creating websearch service: missing FIRECRAWL_API_KEY env")
	}

	return &service{
		as: as,
		c: http.Client{
			Timeout: 1 * time.Minute,
		},
		fireCrawlerApiKey: firecrawlApiKey,
	}, nil
}

type service struct {
	fireCrawlerApiKey string
	as                ai.Service
	c                 http.Client
}

func (s *service) SearchForSpecificPages(ctx context.Context, domains QueryDomains) ([]SearchResult, error) {

	respCh := make(chan SearchResult, len(domains.Queries))
	errCh := make(chan error, len(domains.Queries))
	wg := sync.WaitGroup{}
	for _, query := range domains.Queries {
		wg.Add(1)
		go func(query Query) {
			defer wg.Done()
			prompt := fmt.Sprintf("site: %s, %s", query.Url, query.Q)

			searchRequestBody := FirecrawlSearchRequest{
				Query:   prompt,
				Limit:   3,
				Timeout: 60000,
			}
			req, err := prepareFirecrawlSearchRequest(ctx, s.fireCrawlerApiKey, searchRequestBody)
			if err != nil {
				errCh <- fmt.Errorf("creating request: %w", err)
				return
			}

			resp, err := s.c.Do(req)
			if err != nil {
				errCh <- fmt.Errorf("response error: %w", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				errCh <- fmt.Errorf("response code %d", resp.StatusCode)
				return
			}
			var result FirecrawlSearchResponse

			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				errCh <- fmt.Errorf("decode response: %w", err)
				return
			}

			respCh <- SearchResult{Query: query.Q, Results: result.Results}
		}(query)
	}

	wg.Wait()
	close(errCh)
	close(respCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	var results []SearchResult
	for result := range respCh {
		results = append(results, result)
	}

	return results, nil
}

func (s *service) IsSearchRequired(ctx context.Context, query string) bool {
	msg := []ai.Message{
		ai.SystemMessage(useSearchPrompt),
		ai.UserMessage(query),
	}

	answer, err := s.as.Chat(ctx, msg)
	if err != nil {
		log.Printf("failed to chat with AI: %v", err)
		return false
	}

	return strings.TrimSpace(answer) == doSearch
}

func (s *service) GetDomainQueries(ctx context.Context, query string, allowedDomains []AllowedDomain) (QueryDomains, error) {
	msg := []ai.Message{
		ai.SystemMessage(getAskDomainPrompt(allowedDomains)),
		ai.UserMessage(query),
	}
	answer, err := s.as.Chat(ctx, msg)
	if err != nil {
		log.Printf("failed to get domain queries: %v", err)
		return QueryDomains{}, fmt.Errorf("failed to get domain queries: %w", err)
	}

	return parseDomainQueries(answer)

}

func (s *service) ScoreResults(ctx context.Context, results []SearchResult, userQuery string) ([]WebPage, error) {

	type webPageScore struct {
		webPage WebPage
		scoring float64
	}

	totalPages := 0
	for _, result := range results {
		totalPages += len(result.Results)
	}
	scoringCh := make(chan webPageScore, totalPages)
	errCh := make(chan error, totalPages)
	var wg sync.WaitGroup

	for _, result := range results {
		for _, page := range result.Results {
			wg.Add(1)
			go func(page WebPage, query string) {
				defer wg.Done()
				userPrompt := userScoringPrompt(page, result.Query, userQuery)
				resp, err := s.as.Chat(ctx, []ai.Message{
					ai.SystemMessage(systemScoringPrompt),
					ai.UserMessage(userPrompt),
				})
				if err != nil {
					errCh <- fmt.Errorf("ai response: %w", err)
					return
				}
				var scoringResult scoringAiResponse
				err = json.Unmarshal([]byte(resp), &scoringResult)
				if err != nil {
					errCh <- fmt.Errorf("unmarshal ai response: %w", err)
					return
				}
				scoringCh <- webPageScore{webPage: page, scoring: scoringResult.Score}
			}(page, result.Query)

		}
	}
	wg.Wait()
	close(scoringCh)
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	var webPageScores []webPageScore
	for wps := range scoringCh {
		webPageScores = append(webPageScores, wps)
	}

	sort.Slice(webPageScores, func(i, j int) bool {
		return webPageScores[i].scoring < webPageScores[j].scoring
	})

	var scoredWebPage []WebPage
	for _, v := range webPageScores {
		scoredWebPage = append(scoredWebPage, v.webPage)
	}

	return scoredWebPage[:3], nil

}

func (s *service) ScrapWebpages(ctx context.Context, pages []WebPage, userQuery string) ([]ScrappedWebPage, error) {
	respCh := make(chan ScrappedWebPage, len(pages))
	errCh := make(chan error, len(pages))

	wg := sync.WaitGroup{}
	for _, page := range pages {
		wg.Add(1)
		go func(page WebPage, userQuery string) {
			defer wg.Done()

			firecrawlScrapReq := FirecrawlScrapRequest{
				Url:     page.Url,
				Formats: []FirecrawlScrapFormat{FirecrawlMarkdownFormat},
			}

			req, err := prepareFirecrawlScrapRequest(ctx, s.fireCrawlerApiKey, firecrawlScrapReq)
			if err != nil {
				errCh <- fmt.Errorf("creating request: %w", err)
				return
			}
			resp, err := s.c.Do(req)
			if err != nil {
				errCh <- fmt.Errorf("response error: %w", err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode >= 400 {
				errCh <- fmt.Errorf("response code %d", resp.StatusCode)
				return
			}

			var result FirecrawlScrapResponse
			err = json.NewDecoder(resp.Body).Decode(&result)
			if err != nil {
				errCh <- fmt.Errorf("decode response: %w", err)
				return
			}

			dataMap, ok := result.Data.(map[string]interface{})
			if !ok {
				errCh <- fmt.Errorf("response wront format")
				return
			}
			content, ok := dataMap[string(FirecrawlMarkdownFormat)].(string)

			respCh <- ScrappedWebPage{
				WebPage: page,
				Content: content,
			}
		}(page, userQuery)
	}
	wg.Wait()
	close(respCh)
	close(errCh)

	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	var scrappedWebPages []ScrappedWebPage
	for item := range respCh {
		scrappedWebPages = append(scrappedWebPages, item)
	}
	return scrappedWebPages, nil

}
