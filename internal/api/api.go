package api

import (
	"context"
	"io"
)

// StreamHandler 定义流式传输接口
type StreamHandler interface {
	StreamChat(ctx context.Context, system, user string, writer io.Writer) error
}

// StreamMiddleware 提供流式处理中间件
func StreamMiddleware(next func(context.Context, string, string, io.Writer) error) StreamHandler {
	return &streamHandler{handler: next}
}

type streamHandler struct {
	handler func(context.Context, string, string, io.Writer) error
}

func (h *streamHandler) StreamChat(ctx context.Context, system, user string, writer io.Writer) error {
	return StreamChatCompletion(ctx, system, user, writer)
}