package ai

import (
	"ac/internal/websocket"
	"context"
	"io"
)

type AIClient struct {
	system_prompt string
	broadcaster   websocket.Broadcaster
}

func NewAIClient(system_prompt string, broadcaster websocket.Broadcaster) *AIClient {
	return &AIClient{
		system_prompt: system_prompt,
		broadcaster:   broadcaster,
	}
}

func (c *AIClient) StreamQuery(userMessage string, writer io.Writer) error {
	// 创建自定义writer同时写入原始writer和广播到WebSocket
	broadcastWriter := &broadcastWriter{
		writer:      writer,
		broadcaster: c.broadcaster,
	}

	return StreamChatCompletion(context.Background(), c.system_prompt, userMessage, broadcastWriter)
}

type broadcastWriter struct {
	writer      io.Writer
	broadcaster websocket.Broadcaster
}

func (bw *broadcastWriter) Write(p []byte) (n int, err error) {
	// 写入原始writer
	n, err = bw.writer.Write(p)
	if err != nil {
		return n, err
	}

	// 广播到WebSocket
	bw.broadcaster.Broadcast(websocket.MessageTypeAIResponse, string(p))

	return n, nil
}
