package testdbutils_test

import (
	"testing"

	testdbutils "github.com/djpiper28/rpg-book/desktop_client/test_db_utils"
	"github.com/stretchr/testify/require"
)

func TestDb(t *testing.T) {
	db, close := testdbutils.GetDb()
	defer close()

	require.NotNil(t, db)
}
