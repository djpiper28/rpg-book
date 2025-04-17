package testdbutils_test

import (
	"testing"

	testdbutils "github.com/djpiper28/rpg-book/desktop_client/test_db_utils"
	"github.com/stretchr/testify/require"
)

func TestPrimaryDb(t *testing.T) {
	db, close := testdbutils.GetPrimaryDb()
	defer close()

	require.NotNil(t, db)
}

func TestProjectDb(t *testing.T) {
	db, close := testdbutils.GetProjectDb()
	defer close()

	require.NotNil(t, db)
}
