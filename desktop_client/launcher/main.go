package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/djpiper28/rpg-book/common/platform"
	"github.com/djpiper28/rpg-book/desktop_client/backend"
	"github.com/jessevdk/go-flags"
)

type options struct {
	LogLevel          string `long:"log-level" short:"l" default:"info" description:"Logging level, (info, error, or warning)"`
	WorkdingDirectory string `long:"working-dir" short:"d" description:"Directory to execute the launcher in"`
}

const (
	EnvVarPrefix      = "RPG_BOOK_"
	EnvVarCertificate = EnvVarPrefix + "CERTIFICATE"
	EnvVarPort        = EnvVarPrefix + "PORT"
)

const pathLoggerTag = "path"

func createAppDir(path string) error {
	err := os.MkdirAll(path, 0o755)
	if err != nil {
		return errors.Join(errors.New("Cannot create app data folder"), err)
	}

	log.Info("Created data directory", pathLoggerTag, path)
	return nil
}

func changeWorkingDirectory() (string, error) {
	path := platform.GetAppPath("RPG Book")
	log.Debug("Trying to open data directory", pathLoggerTag, path)

	stat, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err = createAppDir(path)
			if err != nil {
				return "", err
			}

			// If the app directory was created successfully then we can probably set it as our working directory
		} else {
			return "", err
		}
	} else if !stat.IsDir() {
		return "", errors.New("app data folder is not a directory")
	}

	err = os.Chdir(path)
	if err != nil {
		return "", errors.Join(errors.New("Cannot open data directory"), err)
	}

	return path, nil
}

func main() {
	fmt.Println("Starting RPG Book")
	log.Default().SetReportCaller(true)
	log.Default().SetReportTimestamp(true)

	path, err := changeWorkingDirectory()
	if err != nil {
		log.Warn("Cannot set the app to run its data directory", loggertags.TagError, err)
	} else {
		log.Info("Current working directory", pathLoggerTag, path)
	}

	var opts options
	launcherCmd, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("Cannot parse arguments", loggertags.TagError)
	}

	switch strings.ToLower(opts.LogLevel) {
	case "err", "erro", "error":
		log.SetLevel(log.ErrorLevel)
	case "warn", "warning":
		log.SetLevel(log.WarnLevel)
	case "inf", "info":
		log.SetLevel(log.InfoLevel)
	case "dbg", "debug":
		log.SetLevel(log.DebugLevel)
	}

	if len(launcherCmd) == 0 {
		log.Fatal("Using default launcher cmd")
	}

	server, err := backend.New(backend.RandPort())
	if err != nil {
		log.Fatal("Cannot start gRPC server", loggertags.TagError, err)
	}
	defer server.Stop()

	// here to stop certificate leak to logs
	log.Info("Starting RPG book",
		"cmd", launcherCmd,
		"directory", opts.WorkdingDirectory,
		"port", server.Port)

	cmd := exec.Command(launcherCmd[0], launcherCmd[1:]...)
	cmd.Dir = opts.WorkdingDirectory
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, fmt.Sprintf(EnvVarCertificate+"=%s", server.ClientCredentials().CertPem))
	cmd.Env = append(cmd.Env, fmt.Sprintf(EnvVarPort+"=%d", server.ClientCredentials().Port))

	err = cmd.Start()
	if err != nil {
		log.Fatal("Cannot start application", loggertags.TagError, err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Fatal("Cannot wait application", loggertags.TagError, err)
	}

	code := cmd.ProcessState.ExitCode()
	if code != 0 {
		log.Fatal("Application crashed", loggertags.TagExitCode, code)
	}

	log.Infof("Exited without error")
}
