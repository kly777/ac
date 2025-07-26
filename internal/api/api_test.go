package api

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestStreamChat(t *testing.T) {
	// 创建测试用的流式处理器
	handler := StreamMiddleware(func(ctx context.Context, system, user string, writer io.Writer) error {
		// 模拟AI流式响应
		responses := []string{"Hello", " ", "World", "!"}
		for _, resp := range responses {
			if _, err := io.WriteString(writer, resp); err != nil {
				return err
			}
		}
		return nil
	})

	// 创建缓冲区接收流式数据
	var result strings.Builder

	// 执行流式调用
	err := handler.StreamChat(context.Background(), "test-sys", "test-user", &result)
	if err != nil {
		t.Errorf("StreamChat failed: %v", err)
	}

	// 验证结果
	expected := "Hello World!"
	if result.String() != expected {
		t.Errorf("Expected %q, got %q", expected, result.String())
	}
}