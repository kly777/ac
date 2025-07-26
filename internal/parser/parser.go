package parser

import "ac/internal/command"

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Parse(response string) ([]command.Cmd, error) {
	return []command.Cmd{}, nil
}

func (p *Parser) ParseLine(line string) ([]command.Cmd, error) {
	return []command.Cmd{}, nil
}