package websearch

type AllowedDomain struct {
	Domain string `json:"domain"`
	Url    string `json:"url"`
}

type QueryDomains struct {
	Thoughts string  `json:"_thoughts"`
	Queries  []Query `json:"queries"`
}

type Query struct {
	Q   string `json:"q"`
	Url string `json:"url"`
}

type SearchResult struct {
	Query   string    `json:"query"`
	Results []WebPage `json:"results"`
}

type WebPage struct {
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ScrappedWebPage struct {
	WebPage
	Content string `json:"content"`
}

type FirecrawlSearchRequest struct {
	Query         string `json:"query"`
	Limit         int    `json:"limit"`
	Timeout       int    `json:"timeout"`
	ScrapeOptions *struct {
		Formats []string `json:"formats"`
	} `json:"scrapeOptions,omitempty"`
}

type FirecrawlSearchResponse struct {
	Results []WebPage `json:"data"`
}

type FirecrawlScrapFormat string

const (
	FirecrawlMarkdownFormat FirecrawlScrapFormat = "markdown"
	FirecrawlHTMLFormat     FirecrawlScrapFormat = "html"
)

type FirecrawlScrapRequest struct {
	Url     string                 `json:"url"`
	Formats []FirecrawlScrapFormat `json:"formats"`
}

type FirecrawlScrapResponse struct {
	Data interface{} `json:"data"`
}

type scoringAiResponse struct {
	Reason string  `json:"reason"`
	Score  float64 `json:"score"`
}
