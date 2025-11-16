package parser

import (
	"strings"
)

type BinaryOperator int

const (
	BinaryOperator_And = iota + 1
	BinaryOperator_Or
	BinaryOperator_Xor
)

type GeneratorOperator int

const (
	GeneratorOperator_Includes          = iota + 1 // :
	GeneratorOperator_Equals                       // =
	GeneratorOperator_LessThan                     // <
	GeneratorOperator_LessThanEquals               // <=
	GeneratorOperator_GreaterThan                  // >
	GeneratorOperator_GreaterThanEquals            // >=
	GeneratorOperator_NotEqual                     // ~
)

type SetGenerator struct {
	Negated    bool
	Key, Value string
}

type NodeType int

const (
	NodeType_SetGenerator = iota + 1
	NodeType_Basic
	NodeType_Bracket
	NodeType_BinaryOperator
)

type Node struct {
	Left, Right *Node
	Type        NodeType
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
