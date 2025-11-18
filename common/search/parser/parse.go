package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/common/normalisation"
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
	NodeType_Brackets
	NodeType_BinaryOperator
)

type Node struct {
	Left, Right    *Node
	Type           NodeType
	SetGenerator   SetGenerator
	BasicQuery     TextQuery
	BinaryOperator BinaryOperator
}

func Parse(s string) (*Node, error) {
	p := &parser{
		Buffer: strings.ToLower(s),
		root:   &Node{},
		stack:  make([]*Node, 0),
	}
	p.current = p.root
	p.stack = append(p.stack, p.current)

	p.Init()

	err := p.Parse()
	if err != nil {
		log.Error("Failed to parse input", loggertags.TagError, err, "input", s, "tree", p.SprintSyntaxTree())
		return nil, err
	}

	p.Execute()
	err = p.Err
	if err != nil {
		log.Error("Failed to parse input (action failed)", loggertags.TagError, err, "input", s, "tree", p.SprintSyntaxTree())
		return nil, err
	}
	return p.root, nil
}

func NormText(text string) string {
	return normalisation.Normalise(text)
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
