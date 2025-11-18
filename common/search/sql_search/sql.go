package sqlsearch

import (
	"errors"
	"fmt"

	"github.com/djpiper28/rpg-book/common/search/parser"
)

type SqlTableData struct {
	FieldsToScan []string // i.e: []string{"name", "description", ... }; You should use the SqlColmnMap to get the keys of set generators (or the basic query) to map to these
	TableName    string   // i.e: example_table
	JoinClauses  string   // i.e: INNER JOIN Customers ON Orders.CustomerID=Customers.CustomerID;
}

type SqlColmnMap struct {
	TextColumns         map[string]string
	NumberColumns       map[string]string
	BasicQueryColumn    string
	BasicQueryOperation parser.GeneratorOperator // Defaults to includes
}

type sqlScanner struct {
	SqlColmns   SqlColmnMap
	OrderedArgs []any // SqlX takes any for the args
}

func AsSql(query *parser.Node, tableData SqlTableData, columnMap SqlColmnMap) (string, []any, error) {
	// Setup defaults
	if columnMap.BasicQueryOperation == 0 {
		columnMap.BasicQueryOperation = parser.GeneratorOperator_Includes
	}

	if columnMap.NumberColumns == nil {
		columnMap.NumberColumns = make(map[string]string)
	}

	if columnMap.TextColumns == nil {
		columnMap.TextColumns = make(map[string]string)
	}

	// Generate SELECT clause
	selectClause := "SELECT "
	for i, field := range tableData.FieldsToScan {
		selectClause += field
		if i < len(tableData.FieldsToScan)-1 {
			selectClause += ", "
		}
	}

	// Generate FROM caluse
	fromClause := fmt.Sprintf("FROM %s", tableData.TableName)

	if tableData.JoinClauses != "" {
		fromClause += "\n" + tableData.JoinClauses + "\n"
	}

	// Generate WHERE clause
	s := sqlScanner{
		SqlColmns:   columnMap,
		OrderedArgs: make([]any, 0),
	}

	whereClause, err := s.GetQuery(0, query)
	if err != nil {
		return "", nil, errors.Join(errors.New("Cannot create WHERE clause"), err)
	}

	whereClause = "WHERE\n" + whereClause

	return selectClause + "\n" + fromClause + "\n" + whereClause + ";", s.OrderedArgs, nil
}

func nSpaces(n int) string {
	const tabWidth = 2

	output := ""
	for range n * tabWidth {
		output += "  "
	}
	return output
}

const termSpaces = 1

func (s *sqlScanner) ProcessBinaryOperator(depth int, node *parser.Node) (string, error) {
	output := "\n"
	output += nSpaces(depth)
	output += "(\n"

	output += nSpaces(depth + termSpaces)
	if node.Left != nil {
		leftRes, err := s.GetQuery(depth+1, node.Left)
		if err != nil {
			return "", errors.Join(fmt.Errorf("Cannot process left child of node %d (type %v)", depth, node.Type), err)
		}

		output += " " + leftRes
	}

	output += "\n"
	output += nSpaces(depth)

	switch node.BinaryOperator {
	case parser.BinaryOperator_And:
		output += "AND"
	case parser.BinaryOperator_Or:
		output += "OR"
	default:
		return "", fmt.Errorf("%v is not a supported binary operator", node.BinaryOperator)
	}

	output += "\n"
	output += nSpaces(depth + termSpaces)

	if node.Right != nil {
		rightRes, err := s.GetQuery(depth+1, node.Right)
		if err != nil {
			return "", errors.Join(fmt.Errorf("Cannot process right child of node %d (type %v)", depth, node.Type), err)
		}

		output += " " + rightRes
	}

	output += "\n"
	output += nSpaces(depth)
	output += ")"
	return output, nil
}

func (s *sqlScanner) GetQuery(depth int, node *parser.Node) (string, error) {
	switch node.Type {
	case parser.NodeType_Basic:
		return s.ProcessSqlSetGenerator(s.SqlColmns.BasicQueryColumn, s.SqlColmns.BasicQueryOperation, node.BasicQuery.Value)
	case parser.NodeType_SetGenerator:
		return s.ProcessSqlSetGenerator(node.SetGenerator.Key, node.SetGenerator.GeneratorOperator, node.SetGenerator.Value)
	case parser.NodeType_BinaryOperator:
		return s.ProcessBinaryOperator(depth, node)
	default:
		return "", fmt.Errorf("Unsupported node type %v", node.Type)
	}
}

func (s *sqlScanner) ProcessSqlSetGenerator(key string, operator parser.GeneratorOperator, value string) (string, error) {
	if key == "" {
		return "", errors.New("Cannot query with empty keys")
	}

	defer func() {
		s.OrderedArgs = append(s.OrderedArgs, value)
	}()

	sqlOperator := ""
	mappedKey, ok := s.SqlColmns.NumberColumns[key]
	if ok {
		switch operator {
		case parser.GeneratorOperator_GreaterThan:
			sqlOperator = ">"
		case parser.GeneratorOperator_GreaterThanEquals:
			sqlOperator = ">="
		case parser.GeneratorOperator_LessThan:
			sqlOperator = "<"
		case parser.GeneratorOperator_LessThanEquals:
			sqlOperator = "<="
		case parser.GeneratorOperator_Equals:
			sqlOperator = "="
		case parser.GeneratorOperator_NotEquals:
			sqlOperator = "<>"
		default:
			return "", fmt.Errorf("%d is not a supported number generator operator", operator)
		}

		return fmt.Sprintf("%s %s ?", mappedKey, sqlOperator), nil
	}

	mappedKey, ok = s.SqlColmns.TextColumns[key]
	if !ok {
		return "", fmt.Errorf("'%s' is not in the column map", key)
	}

	switch operator {
	case parser.GeneratorOperator_Equals:
		sqlOperator = "="
	case parser.GeneratorOperator_NotEquals:
		sqlOperator = "<>"
	case parser.GeneratorOperator_Includes:
		sqlOperator = "LIKE"
		value = "%" + value + "%"
	default:
		return "", fmt.Errorf("%d is not a supported text generator operator", operator)
	}

	return fmt.Sprintf("%s %s ?", mappedKey, sqlOperator), nil
}
