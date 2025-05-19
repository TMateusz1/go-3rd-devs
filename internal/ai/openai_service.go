package ai

import (
	"context"
	"fmt"
	"github.com/TMateusz1/go-3rd-devs/internal/ai/option"
	"github.com/openai/openai-go"
	option2 "github.com/openai/openai-go/option"
)

type openaiService struct {
	client *openai.Client
	config *option.OpenaiConfig
}

func NewOpenaiService(opts ...option.Option) (Service, error) {
	config, err := option.NewOpenaiConfig(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create openai service: %w", err)
	}
	client := openai.NewClient(option2.WithAPIKey(config.ApiKey), option2.WithBaseURL(config.BaseUrl))

	return &openaiService{
		client: &client,
		config: config,
	}, nil
}

func (o *openaiService) Chat(ctx context.Context, messages []Message) (string, error) {
	return o.ChatWithModel(ctx, messages, o.config.DefaultModel)
}

func (o *openaiService) ChatWithModel(ctx context.Context, messages []Message, model string) (string, error) {
	openaiMessages, err := mapMessagesToOpenaiMessages(messages)
	if err != nil {
		return "", err
	}
	resp, err := o.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    model,
		Messages: openaiMessages,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send message to AI: %w", err)
	}
	return resp.Choices[0].Message.Content, nil
}

func mapMessagesToOpenaiMessages(messages []Message) ([]openai.ChatCompletionMessageParamUnion, error) {
	var result []openai.ChatCompletionMessageParamUnion

	for _, message := range messages {
		openaiMessage, err := mapMessageToOpenaiMessage(message)
		if err != nil {
			return nil, err
		}
		result = append(result, openaiMessage)
	}

	return result, nil
}

func mapMessageToOpenaiMessage(message Message) (openai.ChatCompletionMessageParamUnion, error) {
	switch message.Role {
	case System:
		return openai.SystemMessage(message.Content), nil
	case User:
		return openai.UserMessage(message.Content), nil
	case Assistant:
		return openai.AssistantMessage(message.Content), nil
	default:
		return openai.ChatCompletionMessageParamUnion{}, fmt.Errorf("unknown role: %s", message.Role)
	}
}
