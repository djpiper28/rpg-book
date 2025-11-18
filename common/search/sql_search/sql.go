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

type sqlScanner[T any] struct {
	SqlColmns   SqlColmnMap
	OrderedArgs []string
}

func AsSql[T any](query *parser.Node, tableData SqlTableData, columnMap SqlColmnMap) (string, error) {
	selectClause := "SELECT "
	fromClause := fmt.Sprintf("FROM %s", tableData.TableName)

	if tableData.JoinClauses != "" {
		fromClause += "\n" + tableData.JoinClauses
	}

	// Setup defaults
	if columnMap.BasicQueryOperation == 0 {
		columnMap.BasicQueryOperation = parser.GeneratorOperator_Includes
	}

	// Generate WHERE clause
	s := sqlScanner[T]{
		SqlColmns: columnMap,
	}

	whereClause, err := s.GetQuery(0, query)
	if err != nil {
		return "", errors.Join(errors.New("Cannot create WHERE clause"), err)
	}

	return selectClause + "\n" + fromClause + "\n" + whereClause + ";", nil
}

func (s *sqlScanner[T]) GetQuery(depth int, node *parser.Node) (string, error) {
	switch node.Type {
	case parser.NodeType_Basic:
		return "", errors.ErrUnsupported
	case parser.NodeType_SetGenerator:
		return "", errors.ErrUnsupported
	case parser.NodeType_BinaryOperator:
		output := "\n"
		for range depth {
			output += "  "
		}
		output += "("

		if node.Left != nil {
			leftRes, err := s.GetQuery(depth+1, node.Left)
			if err != nil {
				return "", errors.Join(fmt.Errorf("Cannot process left child of node %d (type %v)", depth, node.Type), err)
			}

			output += " " + leftRes
		}

		switch node.BinaryOperator {
		case parser.BinaryOperator_And:
		case parser.BinaryOperator_Or:
		default:
			return "", fmt.Errorf("%v is not a supported binary operator", node.BinaryOperator)
		}

		if node.Right != nil {
			rightRes, err := s.GetQuery(depth+1, node.Right)
			if err != nil {
				return "", errors.Join(fmt.Errorf("Cannot process right child of node %d (type %v)", depth, node.Type), err)
			}

			output += " " + rightRes
		}

		output += ")\n"
		return output, nil
	default:
		return "", fmt.Errorf("Unsupported node type %v", node.Type)
	}
}

func (s *sqlScanner[T]) ProcessSqlSetGenerator(key string, operator parser.GeneratorOperator, value string) (string, error) {
	sqlOperator := ""
	mappedKey, ok := s.SqlColmns.NumberColumns[key]
	if ok {
		switch operator {
		case parser.GeneratorOperator_GreaterThan:
		case parser.GeneratorOperator_GreaterThanEquals:
		case parser.GeneratorOperator_LessThan:
		case parser.GeneratorOperator_LessThanEquals:
		default:
			return "", fmt.Errorf("%d is not a supported generator operator", operator)
		}

		return fmt.Sprintf("%s%s?", mappedKey, sqlOperator), nil
	}

	mappedKey, ok = s.SqlColmns.TextColumns[key]
	if !ok {
		return "", fmt.Errorf("'%s' is not in the column map", key)
	}

	// TODO: number columns

	switch operator {
	case parser.GeneratorOperator_Equals:
	case parser.GeneratorOperator_NotEquals:
	case parser.GeneratorOperator_Includes:
	default:
		return "", fmt.Errorf("%d is not a supported generator operator", operator)
	}

	return fmt.Sprintf("%s%s?", mappedKey, sqlOperator), nil
}
