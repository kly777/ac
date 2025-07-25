package api

import (
	"ac/pkg/config"
	"net/http"

	openai "github.com/sashabaranov/go-openai"
)

func InitClientConfig() openai.ClientConfig {
	cfg, err := config.Get()
	if err != nil {
		panic(err)
	}
	ClientConfig := openai.DefaultConfig(cfg.DeepSeekAPIKey)
	ClientConfig.BaseURL = "https://api.deepseek.com/v1/"
	ClientConfig.HTTPClient = &http.Client{}

	return ClientConfig
}
