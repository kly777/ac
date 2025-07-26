package api

import (
	"context"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
)

// StreamChatCompletion 流式返回 AI 生成的响应
func StreamChatCompletion(ctx context.Context, systemMessage, userMessage string, writer io.Writer) error {
	client := openai.NewClientWithConfig(InitClientConfig())

	// 创建流式请求
	stream, err := client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: systemMessage},
			{Role: "user", Content: userMessage},
		},
	})
	if err != nil {
		return fmt.Errorf("创建流式请求失败: %w", err)
	}
	defer stream.Close()

	// 实时读取流式数据
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			response, err := stream.Recv()
			if err != nil {
				return fmt.Errorf("接收流式响应失败: %w", err)
			}
			// 将响应内容写入输出（如控制台、WebSocket、HTTP Response）
			if len(response.Choices) > 0 {
				content := response.Choices[0].Delta.Content
				if _, err := io.WriteString(writer, content); err != nil {
					return fmt.Errorf("写入流式数据失败: %w", err)
				}
			}
		}
	}
}