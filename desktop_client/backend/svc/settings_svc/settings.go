package settingssvc

import (
	"context"
	"errors"

	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_common"
	"github.com/djpiper28/rpg-book/desktop_client/backend/pb_settings"
)

type SettingsSvc struct {
	pb_settings.UnimplementedSettingsSvcServer
}

func New() *SettingsSvc {
	return &SettingsSvc{}
}

func (s *SettingsSvc) GetSettings(ctx context.Context, in *pb_common.Empty) (*pb_settings.Settings, error) {
	return nil, errors.New("Not implemented")
}
