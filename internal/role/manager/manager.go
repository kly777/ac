package manager

import (
	"ac/internal/command"
	"ac/internal/info"
	"ac/internal/task"
	"context"
	"log"
)

type AI interface {
	Query(input string) (string, error)
}

type TaskManager interface {
	Add(task task.Task) error
	List() ([]task.Task, error)
}

type Informer interface {
	Get() ([]info.Info, error)
	Format([]info.Info) string
}
type Executor interface {
	Execute(command.Cmd) error
}

type Parser interface {
	Parse(response string) ([]command.Cmd, error)
}

type manager struct {
	taskManager TaskManager
	informer    Informer
	aiClient    AI
	executor    Executor
	parser      Parser
}

func NewManager(tm TaskManager, inf Informer, ai AI, exec Executor, p Parser) *manager {
	return &manager{
		taskManager: tm,
		informer:    inf,
		aiClient:    ai,
		executor:    exec,
		parser:      p,
	}
}

func (m *manager) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 1. 获取用户输入（需实现输入源，如 CLI 或 Web）
			infos, err := m.informer.Get()
			if err != nil {
				log.Printf("收集输入失败: %v", err)
				continue
			}

			input := m.informer.Format(infos)

			// 2. 请求 AI API
			aiResponse, err := m.aiClient.Query(input)
			if err != nil {
				log.Printf("AI 查询失败: %v", err)
				continue
			}

			// 3. 提取指令并执行
			commands, err := m.parser.Parse(aiResponse)

			for _, cmd := range commands {
				err := m.executor.Execute(cmd)
				if err != nil {
					log.Printf("执行%s指令失败: %v", cmd.Command, err)
					continue
				}
			}
		}
	}
}
