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
	"encoding/json"
	"log"
	"net/http"
)

func Run() {
	// 创建WebSocket Hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// 创建广播器
	broadcaster := websocket.NewHubBroadcaster(wsHub)

	// 创建HTTP路由
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头（修复）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		websocket.ServeWS(wsHub, w, r)
	})
	informer := informer.NewInformer(wsHub)
	// 添加接收前端信息的接口
	http.HandleFunc("/q", func(w http.ResponseWriter, r *http.Request) {
		// 设置CORS头（修复）
		log.Println("接收到前端信息")

		origin := r.Header.Get("Origin")

		allowedOrigin := "http://localhost:5173" // 前端开发服务器地址
		if origin == allowedOrigin {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		log.Println("Received request")
		// 处理OPTIONS预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		type InfoRequest struct {
			Content string `json:"content"`
		}

		var req InfoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// 添加到informer
		informer.Add(*info.NewInfo("USER", req.Content))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Info added successfully"))
	})

	// 启动HTTP服务器
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			panic("HTTP server failed: " + err.Error())
		}
	}()

	taskManager := taskManager.NewTaskManager()
	managerAI := ai.NewAIClient(system_prompt, broadcaster)
	managerParser := parser.NewParser()
	executor := executor.NewExecutor(taskManager)
	manager := manager.NewManager(informer, managerAI, executor, managerParser, wsHub)
	ctx := context.Background()
	go func() {
		err := manager.Run(ctx)
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
