package command

type Cmd struct {
	Type cmdType // 命令类型: "RUN" | "TASK"
	Args string  // 命令参数
}

type cmdType string

const (
	CmdTypeRun  cmdType = "RUN"
	CmdTypeTask cmdType = "TASK"
)
