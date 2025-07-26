package parser

import (
	"ac/internal/command"
	"log"
	"strings"
)

const (
	cmdTaskPrefix = "[CMD:TASK]"
	cmdRunPrefix  = "[CMD:RUN]"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseLine(line string) ([]command.Cmd, error) {

	match := strings.Contains(line, "[CMD:")
	if !match {
		return nil, nil // 非指令行直接返回
	}

	cmdContent := strings.TrimSpace(line)

	// 2. 识别指令类型
	switch {
	case strings.HasPrefix(cmdContent, cmdTaskPrefix):
		return []command.Cmd{command.Cmd{
			Type: command.CmdTypeTask,
			Args: strings.TrimPrefix(cmdContent, cmdTaskPrefix),
		}}, nil
	case strings.HasPrefix(cmdContent, cmdRunPrefix):
		return []command.Cmd{command.Cmd{
			Type: command.CmdTypeRun,
			Args: strings.TrimPrefix(cmdContent, cmdRunPrefix),
		}}, nil

	default:
		log.Printf("未知指令类型: %s", cmdContent)
		return nil, nil
	}
}
