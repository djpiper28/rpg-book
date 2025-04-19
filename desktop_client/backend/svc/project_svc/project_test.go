package projectsvc_test

import (
	"testing"

	testdbutils "github.com/djpiper28/rpg-book/desktop_client/test_db_utils"
)

func TestProjectSvc(t *testing.T) {
	_, closeDb := testdbutils.GetPrimaryDb()
	defer closeDb()
}
