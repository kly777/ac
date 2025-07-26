package api

import (
	"context"
	"io"
)

type AIClient struct {
	system_prompt string
}

func NewAIClient(system_prompt string) *AIClient {
	return &AIClient{
		system_prompt: system_prompt,
	}
}

func (c *AIClient) StreamQuery(userMessage string, writer io.Writer) error {
	return StreamChatCompletion(context.Background(), c.system_prompt, userMessage, writer)
}
