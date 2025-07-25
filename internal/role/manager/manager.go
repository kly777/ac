package manager

import (
	"ac/internal/api"
	"ac/internal/task"
	"ac/internal/db"
	"ac/internal/executor"
	"encoding/json"
)

const system_prompt = `
你是本编程项目的管理者，负责将用户需求拆解为可并行开发的独立任务。请按以下规则生成JSON指令：

1. 通用规范：
   - 使用标准JSON格式
   - 包含版本号（version）和目标语言（target_language）
   - 每个任务必须包含唯一ID、文件路径、详细任务描述

2. 导出规范（exports）：
   - 使用统一元信息结构描述导出元素：
     {
       "type": "struct/function/interface/variable/class/other",
       "name": "元素名称",
       "parameters": [参数列表],
       "returns": [返回值列表],
       "description": "功能描述(详细到如何使用)"
     }
   - 支持多语言特性：
     * Go: 包路径、接口方法签名
     * Python: 类继承关系、装饰器
     * Java: 泛型参数、异常声明

3. 依赖规范（dependencies）：
   - 每个依赖项必须包含：
     {
       "file": "依赖文件路径",
       "required_exports": ["所需导出元素名称"],
       "version": "兼容版本范围"
     }

4. 示例输出：
{
  "version": "1.1",
  "target_language": "go",
  "tasks": [
    {
      "id": "task-001",
      "file": "internal/task/model",
      "task": "实现任务数据结构，包含ID、描述、状态字段，提供NewTask函数",
      "exports": [
        {
          "type": "struct",
          "name": "Task",
          "fields": {
            "ID": "string",
            "Description": "string",
            "Status": "string"
          },
          "access": "public",
          "description": "任务基础数据结构"
        },
        {
          "type": "function",
          "name": "NewTask",
          "parameters": ["description string"],
          "returns": ["*Task"],
          "access": "public",
          "description": "创建新任务实例"
        }
      ],
      "dependencies": []
    },
    {
      "id": "task-002",
      "file": "internal/task/storage",
      "task": "实现任务持久化存储接口，包含Create、Get、List方法",
      "exports": [
        {
          "type": "interface",
          "name": "Storage",
          "methods": {
            "Create": {
              "params": ["task *Task"],
              "returns": ["error"]
            },
            "Get": {
              "params": ["id string"],
              "returns": ["*Task", "error"]
            }
          },
          "description": "任务存储接口"
        }
      ],
      "dependencies": [
        {
          "file": "internal/task/model",
          "required_exports": ["Task", "NewTask"],
          "version": "^1.0"
        }
      ]
    }
  ]
}
`

func A() {
	jsonStr, err := api.Answer(system_prompt, "请用go实现一个任务管理系统，支持创建带依赖关系的任务，并将文档存储到SQLite。每个任务包含ID、描述、依赖关系和文档内容。")
	if err != nil {
		panic(err)
	}

	// 解析JSON响应
	var result struct {
		Version    string          `json:"version"`
		Language   string        `json:"target_language"`
		Tasks      []task.Task   `json:"tasks"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		panic(err)
	}

	// 初始化数据库
	if err := db.InitDB(); err != nil {
		panic(err)
	}

	// 创建任务管理器
	taskManager := NewTaskManager()

	// 创建所有任务
	for _, t := range result.Tasks {
		if err := taskManager.CreateTask(&t); err != nil {
			panic(err)
		}
	}

	// 执行任务
	executor := executor.NewTaskExecutor()
	if err := executor.ExecuteTasks(); err != nil {
		panic(err)
	}
}
