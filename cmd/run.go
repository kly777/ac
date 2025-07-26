package main

import (
	"ac/internal/ai"
	"ac/internal/executor"
	"ac/internal/info"
	"ac/internal/informer"
	"ac/internal/parser"
	"ac/internal/role/manager"
	"ac/internal/taskManager"
	"ac/internal/websocket"
	"context"
	"net/http"
)

func Run() {
	// 创建WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// 创建HTTP路由
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWS(wsHub, w, r)
	})
	
	// 启动HTTP服务器
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic("WebSocket server failed: " + err.Error())
		}
	}()

	// 原有业务逻辑
	informer := informer.NewInformer(wsHub)
	informer.Add(*info.NewInfo("用户输入", "用go写一个控制台累加器"))
	taskManager := taskManager.NewTaskManager()
	managerAI := ai.NewAIClient(system_prompt, wsHub)
	managerParser := parser.NewParser()
	executor := executor.NewExecutor(taskManager)
	manager := manager.NewManager(informer, managerAI, executor, managerParser, wsHub)
	ctx := context.Background()
	go func(){
		err:= manager.Run(ctx)
		if err != nil {
			panic(err)
		}
	}()
	select {} // 阻塞主 goroutine，保持程序运行
}

const system_prompt = `
你是本编程项目的管理者，负责将用户需求拆解为可并行开发的独立任务。你不仅可以输出文字和思考,还可以在输出中采用命令

要考虑整体的架构设计
要考虑任务的分配合理性

## 命令列表
### RUN <command>
	在控制台执行命令
### TASK <json>
	创建一个任务，json格式包含任务名称、描述和依赖
	会传递给其他ai执行的可并行开发的独立任务。

请严格按照以下规则输出指令：

- 所有指令必须用<<<\n与\n>>>包裹。
- 每个指令类型用 [CMD] 标记，格式为 [CMD:TYPE]。
- 支持的指令类型：
  - [CMD:TASK] 创建任务
  - [CMD:RUN] 执行终端命令

示例格式：
[CMD:TASK] {"id": "task-001", "file": "internal/task/model", "task": "实现任务数据结构...", "dependencies": []}

[CMD:RUN] go run main.go
`
