package ai

import (
	"context"
	"io"
)

// StreamHandler 定义流式处理接口，用于处理需要实时传输的会话（如AI对话）。
// 方法约定：
// - StreamChat 将用户输入和系统指令传递给处理函数，并通过 writer 实时返回响应。
type StreamHandler interface {
	StreamChat(ctx context.Context, system, user string, writer io.Writer) error
}

// StreamMiddleware 提供流式处理中间件功能，用于封装通用处理逻辑（如日志、鉴权）。
// 参数 next 是核心处理函数，接收上下文、系统指令、用户输入和响应写入器。
// 返回值为封装后的 StreamHandler 实例。
func StreamMiddleware(next func(context.Context, string, string, io.Writer) error) StreamHandler {
	return &streamHandler{handler: next}
}

// streamHandler 是 StreamHandler 接口的具体实现，包含实际处理逻辑。
type streamHandler struct {
	handler func(context.Context, string, string, io.Writer) error // 封装的核心处理函数
}

// StreamChat 实现流式对话处理，直接调用包装的处理函数。
// 参数说明：
// - ctx 上下文用于控制请求生命周期
// - system 系统指令（如角色定义）
// - user 用户输入内容
// - writer 用于实时写入响应数据流
// 返回值：处理过程中发生的错误
func (h *streamHandler) StreamChat(ctx context.Context, system, user string, writer io.Writer) error {
	// 直接调用包装的处理函数，确保测试时可绕过中间件链直接模拟核心逻辑
	return h.handler(ctx, system, user, writer)
}
