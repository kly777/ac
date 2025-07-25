package api

import (
	"context"
	openai "github.com/sashabaranov/go-openai"
)

func Conv(systemMessage, userMessage string) error {

	client := openai.NewClientWithConfig(InitClientConfig())

	chatCompletion, err := client.CreateChatCompletion(context.TODO(), openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: systemMessage},
			{Role: "user", Content: userMessage},
		},
	})
	if err != nil {
		panic(err.Error())
	}
	println(chatCompletion.Choices[0].Message.Content)
	return nil
}
