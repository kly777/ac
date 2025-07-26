package parser

type AIParser struct{}
func (p *AIParser) Parse(response string) ([]Command, error) {
    // 实现解析规则（如正则匹配、JSON 解析）
    // 示例：提取 "CREATE TASK: ..." 或 "RUN CMD: ..." 指令
    var commands []Command
    // ... 解析逻辑 ...
    return commands, nil
}