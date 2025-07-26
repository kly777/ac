package manager

import (
	"ac/internal/command"
	"ac/internal/info"
	"bufio"
	"context"
	"io"
	"log"
	"strings"
	"time"
)

// AI 接口，用于与 AI 服务交互，获取响应
type ai interface {
	StreamQuery(input string, writer io.Writer) error // 流式输出直接写入 writer
}

// 用户的输入,命令的输出,task完成情况都会加入informer,Manager从informer中获取信息
type informer interface {
	Get() []info.Info
	Clear() // 清空当前的Infos
	Format([]info.Info) string
}

// 执行器接口，用于执行解析后的命令,命令可以创建task,执行终端命令
type executor interface {
	Execute(command.Cmd) error
}

// 解析器接口，用于解析AI的响应，提取出可执行的命令
type parser interface {
	ParseLine(line string) ([]command.Cmd, error) // 解析单行
}

type manager struct {
	informer informer
	aiClient ai
	executor executor
	parser   parser

	streamReader *bufio.Reader  // 流式数据读取器
	streamWriter *io.PipeWriter // 流式数据写入端
}

func NewManager(inf informer, ai ai, exec executor, p parser) *manager {
	// 创建管道：写入端 pw 用于接收流式输出，读取端 pr 用于实时解析
	pr, pw := io.Pipe()
	return &manager{
		informer:     inf,
		aiClient:     ai,
		executor:     exec,
		parser:       p,
		streamReader: bufio.NewReader(pr),
		streamWriter: pw,
	}
}

func (m *manager) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// 1. 获取用户输入
			infos := m.informer.Get()
			log.Println("获取到的用户输入:", infos)
			input := m.informer.Format(infos)
			m.informer.Clear()

			if len(input) == 0 {
				log.Println("没有用户输入，等待下一次循环")
				time.Sleep(time.Second)
				continue
			}
			log.Println("格式化后的输入:", input)

			pr, pw := io.Pipe()
			m.streamReader = bufio.NewReader(pr)
			m.streamWriter = pw

			// 2. 启动解析协程（实时解析流式输出）
			go m.startStreamingParser()

			// 3. 执行流式请求（数据写入管道）
			log.Println("开始执行流式请求")
			err := m.aiClient.StreamQuery(input, m.streamWriter)
			if err != nil {
				m.streamWriter.CloseWithError(err) // 通知解析协程异常
				log.Printf("AI 流式请求失败: %v", err)
				continue
			}
			log.Println("AI 流式请求完成")
		}
	}
}

func (m *manager) startStreamingParser() {
	for {
		line, err := m.streamReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("流式解析完成")
			} else {
				log.Printf("解析流式数据时出错: %v", err)
			}
			break
		}

		// 解析单行内容
		if strings.TrimSpace(line) == "" {
			continue
		}
		log.Println("输出:", line)
		cmds, err := m.parser.ParseLine(line)
		if err != nil {
			log.Printf("解析指令失败: %v", err)
			continue
		}
		if len(cmds) != 0 {
			log.Println("解析到指令:", cmds)
		}

		// 执行命令
		for _, cmd := range cmds {
			if err := m.executor.Execute(cmd); err != nil {
				log.Printf("执行指令失败: %v", err)
			}
		}
	}
}
