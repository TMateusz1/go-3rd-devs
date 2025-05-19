package option

import (
	"fmt"
	"os"
)

const (
	defaultModel   = "gpt-3.5-turbo"
	defaultBaseUrl = "https://api.openai.com/v1"
)

type OpenaiConfig struct {
	ApiKey       string
	DefaultModel string
	BaseUrl      string
}

type Option func(*OpenaiConfig)

func NewOpenaiConfig(opts ...Option) (*OpenaiConfig, error) {
	config := &OpenaiConfig{
		DefaultModel: defaultModel,
		BaseUrl:      defaultBaseUrl,
	}

	loadFromEnv(config)

	for _, opt := range opts {
		opt(config)
	}

	err := validate(config)
	if err != nil {
		return nil, fmt.Errorf("openai config is invalid: %w", err)
	}

	return config, nil
}

func validate(config *OpenaiConfig) error {
	// should have more real validators like Parsing URL, model validation etc.
	if config.ApiKey == "" {
		return fmt.Errorf("apikey is required")
	}
	if config.DefaultModel == "" {
		return fmt.Errorf("model is empty")
	}
	if config.BaseUrl == "" {
		return fmt.Errorf("baseurl is empty")
	}
	return nil
}

func loadFromEnv(config *OpenaiConfig) {
	apikey, ok := os.LookupEnv("OPENAI_API_KEY")
	if ok {
		config.ApiKey = apikey
	}

	model, ok := os.LookupEnv("OPENAI_MODEL")
	if ok {
		config.DefaultModel = model
	}

	baseUrl, ok := os.LookupEnv("OPENAI_BASE_URL")
	if ok {
		config.BaseUrl = baseUrl
	}

}

func WithApiKey(apiKey string) Option {
	return func(config *OpenaiConfig) {
		config.ApiKey = apiKey
	}
}

func WithBaseModel(model string) Option {
	return func(config *OpenaiConfig) {
		config.DefaultModel = model
	}
}

func WithBaseUrl(baseUrl string) Option {
	return func(config *OpenaiConfig) {
		config.BaseUrl = baseUrl
	}
}
