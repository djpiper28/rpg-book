package platform

import (
	"os"
	"path"
	"runtime"

	"github.com/charmbracelet/log"
)

func GetAppPath(appname string) string {
	basePath := ""

	switch runtime.GOOS {
	case "windows":
		home := os.Getenv("HOME")
		basePath = path.Join(home, "AppData", "Local")
	case "linux", "freebsd":
		home := os.Getenv("XDG_DATA_HOME")
		basePath = home
	case "darwin":
		home := os.Getenv("HOME")
		basePath = path.Join(home, "Library")
	default:
		log.Warn("Unknown OS, defaulting to current directory. This OS might work, but who knows.", "os", runtime.GOOS)
	}

	if basePath == "" {
		basePath = "."
	}

	basePath = path.Join(basePath, appname)
	return basePath
}
