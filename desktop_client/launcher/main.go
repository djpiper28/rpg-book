package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	LogLevel          string `long:"log" short:"l" default:"info" description:"Logging level, (info, error, or warning)"`
	WorkdingDirectory string `long:"dir" short:"d" description:"Directory to execute the launcher in"`
}

func main() {
	fmt.Println("Starting RPG Book")
	log.Default().SetReportCaller(true)
	log.Default().SetReportTimestamp(true)

	var opts Options
	launcherCmd, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal("Cannot parse arguments", loggertags.TagError)
	}

	switch strings.ToLower(opts.LogLevel) {
	case "err":
		fallthrough
	case "erro":
		fallthrough
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		fallthrough
	case "warning":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	}

	if len(launcherCmd) == 0 {
		log.Info("Using default launcher cmd")
		panic("TODO: set default launcher cmd")
	}

	// TODO: start a server, and append the server port to the arguments

	log.Infof("Starting RPG book with %+v, directory %s", launcherCmd, opts.WorkdingDirectory)
	cmd := exec.Command(launcherCmd[0], launcherCmd[1:]...)
	cmd.Dir = opts.WorkdingDirectory
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

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
}
