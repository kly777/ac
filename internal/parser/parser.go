package parser

import "context"

type Parser interface {
	Parse(response string) ([]Command, error)
}

type Command interface {
	Execute(ctx context.Context) (Cmd string, Result string, err error)
}
