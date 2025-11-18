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
	NodeType_BinaryOperator
)

type Node struct {
	Left, Right    *Node
	Type           NodeType
	SetGenerator   SetGenerator
	BasicQuery     TextQuery
	BinaryOperator BinaryOperator
}

func (p *parser) insertOperatorNode() {
	newNode := &Node{
		Type:           NodeType_BinaryOperator,
		BinaryOperator: p.tempNode.BinaryOperator,
	}

	p.stack = append(p.stack, newNode)
}

func (p *parser) panic(reason error) {
	log.Error("Panic whilst parsing input",
		loggertags.TagError, reason,
		"tree", p.SprintSyntaxTree(),
		"stack", p.stack,
		"tempNode", p.tempNode)
	panic(reason)
}

func (p *parser) pop() *Node {
	if len(p.stack) == 0 {
		p.panic(errors.New("Stack is empty so cannot be popped"))
	}

	top := p.stack[len(p.stack)-1]
	p.stack = p.stack[:len(p.stack)-1]
	return top
}

func (p *parser) push(node *Node) {
	p.stack = append(p.stack, node)
	p.tempNode = &Node{}
}

func (p *parser) finaliseStack() *Node {
	if len(p.stack) == 1 {
		top := p.stack[0]
		p.stack = make([]*Node, 0)
		return top
	}

	operandA := p.pop()
	operator := p.pop()
	operandB := p.pop()

	if operator.Type != NodeType_BinaryOperator {
		p.panic(fmt.Errorf("Expected an operator, found %+v", operator))
	}

	if operator.Left != nil || operator.Right != nil {
		p.panic(fmt.Errorf("Expected nil children, found %+v", operator))
	}

	operator.Left = operandA
	operator.Right = operandB
	return operator
}

func Parse(inputStr string) (*Node, error) {
	p := &parser{
		Buffer:   strings.ToLower(inputStr),
		tempNode: &Node{},
		stack:    make([]*Node, 0),
	}

	p.Init()

	err := p.Parse()
	if err != nil {
		log.Error("Failed to parse input",
			loggertags.TagError, err,
			"input", inputStr,
			"tree", p.SprintSyntaxTree())
		return nil, err
	}

	p.Execute()
	err = p.Err
	if err != nil {
		log.Error("Failed to parse input (action failed)",
			loggertags.TagError, err,
			"input", inputStr,
			"tree", p.SprintSyntaxTree())
		return nil, err
	}

	top := p.finaliseStack()
	if len(p.stack) != 0 {
		err = fmt.Errorf("Expected the stack to be empty after parsing, it has len %d", len(p.stack))
		log.Error("Failed to parse input (AST creation failed)",
			loggertags.TagError, err,
			"input", inputStr,
			"tree", p.SprintSyntaxTree(),
			"stack", p.stack)
		return nil, err
	}
	return top, nil
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
