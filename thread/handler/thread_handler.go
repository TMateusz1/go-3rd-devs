package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TMateusz1/go-3rd-devs/internal/ai"
	"log"
	"net/http"
)

type ThreadHandler struct {
	as      ai.Service
	summary string
}

func NewThreadHandler(as ai.Service) *ThreadHandler {
	return &ThreadHandler{
		as: as,
	}
}

func (h *ThreadHandler) Handle(w http.ResponseWriter, r *http.Request) {
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

	questionMessages := []ai.Message{
		ai.SystemMessage(fmt.Sprintf("You are a helpful assistant who speaks using as few words as possible. <summary>%s</summary>", h.summary)),
		ai.UserMessage(req.Message),
	}

	answer, err := h.as.Chat(r.Context(), questionMessages)
	if err != nil {
		log.Printf("failed to chat with AI: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.setNewSummary(r.Context(), req.Message, answer)
	if err != nil {
		log.Printf("failed to set new summary: %v", err)
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

func (h *ThreadHandler) setNewSummary(ctx context.Context, message string, answer string) error {
	newSummary, err := h.as.Chat(ctx, []ai.Message{
		ai.SystemMessage(getNewSummaryPrompt(h.summary, message, answer)),
		ai.UserMessage("Please summarize conversation in short way."),
	})
	if err != nil {
		return fmt.Errorf("failed to chat with AI: %w", err)
	}
	h.summary = newSummary
	return nil
}

func getNewSummaryPrompt(previousSummary, userMessage, assistantMessage string) string {
	return fmt.Sprintf(`
Please summarize the following conversation in a concise manner, incorporating the previous summary if available:
<previous_summary>%s</previous_summary>
<current_turn> 
User: %s 
Assistant: %s 
</current_turn>
`, previousSummary, userMessage, assistantMessage)
}
