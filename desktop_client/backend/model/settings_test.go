package model_test

import (
	"testing"

	"github.com/djpiper28/rpg-book/desktop_client/backend/model"
	"github.com/stretchr/testify/require"
)

func TestToAndFromProto(t *testing.T) {
	s := model.Settings{
		DevMode:  true,
		DarkMode: false,
	}

	s2 := model.SettingsFromProto(s.ToProto())
	require.Equal(t, s, *s2)
}

func TestToAndFromProto2(t *testing.T) {
	s := model.Settings{
		DevMode:  false,
		DarkMode: true,
	}

	s2 := model.SettingsFromProto(s.ToProto())
	require.Equal(t, s, *s2)
}
