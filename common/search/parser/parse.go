package parser

import (
	"strings"
)

func Parse(s string) error {
	p := &parser{
		Buffer: strings.ToLower(s),
	}

	p.Init()

	if err := p.Parse(); err != nil {
		return err
	}

	p.Parse()
	// p.PrintSyntaxTree()
	return nil
}
