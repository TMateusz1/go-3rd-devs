package ai

import (
	"context"
)

const (
	OpenAiGPT35Turbo             string = "gpt-3.5-turbo"
	OpenAiGPT4oMini              string = "gpt-4o-mini"
	OpenRouterModelGPT4oMini     string = "openai/gpt-4o-mini"
	OpenRouterModelGPT41Nano     string = "openai/gpt-4.1-nano"
	OpenRouterModelGemini25Flash string = "google/gemini-2.5-flash-preview"
)

type Service interface {
	Chat(ctx context.Context, messages []Message) (string, error)
	ChatWithModel(ctx context.Context, messages []Message, model string) (string, error)
}
