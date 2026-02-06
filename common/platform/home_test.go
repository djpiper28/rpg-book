package platform_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/common/platform"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestHome(t *testing.T) {
	appname := uuid.New().String()
	apppath := platform.GetAppPath(appname)

	t.Log(apppath)
	require.NotEmpty(t, apppath)
	require.Contains(t, apppath, appname)
}
