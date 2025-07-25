package api

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

// StreamAnswer 使用官方openai-go调用DeepSeek API进行流式回答
func Answer(systemMessage, userMessage string) (string, error) {
	// 创建DeepSeek兼容客户端
	client := openai.NewClientWithConfig(InitClientConfig())

	chatCompletion, err := client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: systemMessage},
			{Role: "user", Content: userMessage},
		},
		ResponseFormat: &openai.ChatCompletionResponseFormat{Type: "json_object"},
	})
	if err != nil {
		panic(err.Error())
	}
	v := chatCompletion.Choices[0].Message.Content
	fmt.Println(v)
	return v, nil
}
