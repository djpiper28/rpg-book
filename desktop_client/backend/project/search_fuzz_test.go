package project_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/djpiper28/rpg-book/common/search/parser"
	testutils "github.com/djpiper28/rpg-book/common/test_utils"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func FuzzSearchCharacter(f *testing.F) {
	project, remove := testutils.NewProject(f)
	defer remove()

	for i := range 100 {
		name := fmt.Sprintf("character-%d", i)
		desc := uuid.New().String()
		_, err := project.CreateCharacter(name, desc, nil)
		require.NoError(f, err)

		f.Add(name)
		f.Add(desc)
		f.Add(fmt.Sprintf(`"name":"%s"`, name))
		f.Add(fmt.Sprintf("name:%s", name))
		f.Add(fmt.Sprintf("desc:%s", desc))
		f.Add(fmt.Sprintf("description:%s", desc))

		f.Add(fmt.Sprintf("name=%s", name))
		f.Add(fmt.Sprintf("desc=%s", desc))
	}

	f.Add("name>5")
	f.Add("desc<5")
	f.Add("name>=5")
	f.Add("desc<=5")
	f.Add(`name:"dave"`)
	f.Add(`description:dave`)
	f.Add(`name:"--'; DROP TABLE characters CASCADE; --`)

	f.Fuzz(func(t *testing.T, query string) {
		_, parseErr := parser.Parse(query)
		_, err := project.SearchCharacter(query)

		if parseErr != nil {
			require.Error(t, err)
		} else if err != nil {
			if strings.Contains(err.Error(), "Cannot create WHERE clause") {
				return
			}

			require.NoError(t, err)
		}
	})
}
