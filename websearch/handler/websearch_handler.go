package handler

import (
	"encoding/json"
	"fmt"
	"github.com/TMateusz1/go-3rd-devs/internal/ai"
	"github.com/TMateusz1/go-3rd-devs/internal/websearch"
	"log"
	"net/http"
	"strings"
)

var allowedDomains = []websearch.AllowedDomain{
	{"Wikipedia.org", "https://en.wikipedia.org"},
	{"OpenAI", "https://openai.com"},
	{"Go DEV", "https://go.dev"},
	{"Ardan Labs Golang courses!", "https://www.ardanlabs.com"},
}

type WebSearchHandler struct {
	as ai.Service
	ws websearch.Service
}

func NewWebSearchHandler(as ai.Service, ws websearch.Service) *WebSearchHandler {
	return &WebSearchHandler{
		as: as,
		ws: ws,
	}
}

func (h *WebSearchHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type request struct {
		Message string `json:"message"`
	}

	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var scrappedWebPage []websearch.ScrappedWebPage

	required := h.ws.IsSearchRequired(r.Context(), req.Message)
	if required {
		queries, err := h.ws.GetDomainQueries(r.Context(), req.Message, allowedDomains)
		if err != nil {
			log.Printf("failed to get domain queries: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		results, err := h.ws.SearchForSpecificPages(r.Context(), queries)
		if err != nil {
			log.Printf("error searching for specific domains: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		scoredResults, err := h.ws.ScoreResults(r.Context(), results, req.Message)
		if err != nil {
			log.Printf("error scoring: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		scrappedWebPage, err = h.ws.ScrapWebpages(r.Context(), scoredResults, req.Message)
		if err != nil {
			log.Printf("error scrapping: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	answer, err := h.as.Chat(r.Context(), []ai.Message{
		ai.SystemMessage(promptWithResults(scrappedWebPage)),
		ai.UserMessage(req.Message),
	})

	if err != nil {
		log.Printf("failed to get domain queries: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(struct {
		Answer string `json:"answer"`
	}{answer})

	if err != nil {
		log.Printf("failed to encode response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func promptWithResults(pages []websearch.ScrappedWebPage) string {
	builder := strings.Builder{}
	builder.WriteString("Answer the question based on")
	if pages == nil {
		builder.WriteString("your existing knowledge.\n")
	} else {
		builder.WriteString("provided search results and scraped content.\n")

		builder.WriteString("<search_results>\n")
		for _, page := range pages {
			builder.WriteString(fmt.Sprintf("<search_result url=%s title=%s description=%s>\n", page.Url, page.Title, page.Description))
			builder.WriteString(page.Content)
			builder.WriteString("</search_result>\n")
		}
		builder.WriteString("</search_results>\n")

	}
	builder.WriteString("Use the fewest words possible.")
	return builder.String()
}
