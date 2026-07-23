package parser_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	"github.com/stretchr/testify/require"
)

func FuzzParse(f *testing.F) {
	f.Add("Dave")
	f.Add("John Major")
	f.Add("deleted=1")
	f.Add("deleted=t")
	f.Add("deleted=T")
	f.Add("deleted=true")
	f.Add("deleted=TrUe")

	f.Add("deleted=0")
	f.Add("deleted=f")
	f.Add("deleted=F")
	f.Add("deleted=false")
	f.Add("deleted=FaLsE")

	f.Add("rating>123")
	f.Add("rating>=124")
	f.Add("rating<1000")
	f.Add("rating<=999")
	f.Add("rating=999")
	f.Add("rating~999")

	f.Add("deviation=123")
	f.Add("name=John Major")
	f.Add("name~John Major")
	f.Add(`name~"John Major"`)

	f.Add(`deleted=false and (name_norm=greg OR name_norm="chas") and (rating>600 and rating<1600)`)
	f.Add(`(name_norm=greg OR name_norm="chas") and rating>600 and rating<1600`)

	// Here are some of the fun cases fuzzing raised that have since been fixed
	f.Add("(((((((((((((((((( 0")

	f.Fuzz(func(t *testing.T, query string) {
		t.Logf("Query is: %s", query)
		// This only looks for things that hang, or crash.
		ast, err := parser.Parse(query)
		if err == nil {
			require.NotEmpty(t, ast)
		}
	})
}
