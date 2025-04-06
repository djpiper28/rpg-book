package settingssvc_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_settings"
	settingssvc "github.com/djpiper28/rpg-book/desktop_client/backend/svc/settings_svc"
	"github.com/stretchr/testify/require"
)

var settings pb_settings.SettingsSvcServer = settingssvc.New()

func TestNew(t *testing.T) {
	require.NotNil(t, settings)
}
