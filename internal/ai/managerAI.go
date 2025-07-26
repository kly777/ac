package ai

import (
	"context"
	"io"
	"ac/internal/websocket"
	"encoding/json"
)

type AIClient struct {
	system_prompt string
	hub           *websocket.Hub
}

func NewAIClient(system_prompt string, hub *websocket.Hub) *AIClient {
	return &AIClient{
		system_prompt: system_prompt,
		hub:           hub,
	}
}

func (c *AIClient) StreamQuery(userMessage string, writer io.Writer) error {
	// 创建自定义writer同时写入原始writer和广播到WebSocket
	broadcastWriter := &broadcastWriter{
		writer: writer,
		hub:    c.hub,
	}
	
	return StreamChatCompletion(context.Background(), c.system_prompt, userMessage, broadcastWriter)
}

type broadcastWriter struct {
	writer io.Writer
	hub    *websocket.Hub
}

func (bw *broadcastWriter) Write(p []byte) (n int, err error) {
	// 写入原始writer
	n, err = bw.writer.Write(p)
	if err != nil {
		return n, err
	}
	
	// 广播到WebSocket
	jsonData, _ := json.Marshal(map[string]interface{}{
		"type": "managerAI",
		"data": string(p),
	})
	bw.hub.Broadcast(jsonData)
	
	return n, nil
}
