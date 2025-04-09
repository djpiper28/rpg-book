package migrations

import "testing"

type Migration struct {
	Sql  string
	Test func(testing.T)
}
