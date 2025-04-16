package model

import "github.com/djpiper28/rpg-book/desktop_client/backend/pb_system"

type Settings struct {
	DevMode  bool `db:"dev_mode"`
	DarkMode bool `db:"dark_mode"`
}

func (s *Settings) ToProto() *pb_system.Settings {
	return &pb_system.Settings{
		DevMode:  s.DevMode,
		DarkMode: s.DarkMode,
	}
}

func SettingsFromProto(p *pb_system.Settings) *Settings {
	return &Settings{
		DevMode:  p.DevMode,
		DarkMode: p.DarkMode,
	}
}
