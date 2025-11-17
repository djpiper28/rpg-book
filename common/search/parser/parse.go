package parser

import (
	"errors"
	"fmt"
	"strings"
)

type BinaryOperator int

const (
	BinaryOperator_And BinaryOperator = iota + 1
	BinaryOperator_Or
	BinaryOperator_Xor
	Default_BinaryOperator = BinaryOperator_And
)

type GeneratorOperator int

const (
	GeneratorOperator_Includes          GeneratorOperator = iota + 1 // :
	GeneratorOperator_Equals                                         // =
	GeneratorOperator_LessThan                                       // <
	GeneratorOperator_LessThanEquals                                 // <=
	GeneratorOperator_GreaterThan                                    // >
	GeneratorOperator_GreaterThanEquals                              // >=
	GeneratorOperator_NotEquals                                      // ~
)

type SetGenerator struct {
	Negated           bool
	Key, Value        string
	GeneratorOperator GeneratorOperator
}

type TextQuery struct {
	Value string
}

type NodeType int

const (
	NodeType_SetGenerator NodeType = iota + 1
	NodeType_Basic
	NodeType_Bracket
	NodeType_BinaryOperator
)

type Node struct {
	Left, Right    *Node
	Type           NodeType
	SetGenerator   SetGenerator
	TextQuery      TextQuery
	BinaryOperator BinaryOperator
}

func Parse(s string) error {
	p := &parser{
		Buffer: strings.ToLower(s),
		root:   &Node{},
		stack:  make([]*Node, 0),
	}
	p.current = p.root
	p.stack = append(p.stack, p.current)

	p.Init()

	if err := p.Parse(); err != nil {
		return err
	}

	p.Parse()
	// p.PrintSyntaxTree()
	return nil
}

func UnescapeText(text string) (string, error) {
	var output string
	escaping := false

	for _, c := range text {
		if escaping {
			escaping = false

			switch c {
			case '\\':
				output += "\\"
			case 'r':
				output += "\r"
			case 'n':
				output += "\n"
			case 't':
				output += "\t"
			case '"':
				output += `"`
			case '\'':
				output += "'"
			default:
				return "", fmt.Errorf("'%c' is an invalid escape character", c)
			}
		} else if c == '\\' {
			escaping = true
		} else {
			output += string(c)
		}
	}

	if escaping {
		return "", errors.New("Cannot escape next character due to end of input")
	}

	return output, nil
}
