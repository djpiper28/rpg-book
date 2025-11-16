package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
)

func readOutput(cmd *exec.Cmd) string {
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Cannot create pipe", loggertags.TagError, err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal("Cannot start app", loggertags.TagError, err)
	}

	bytes, err := io.ReadAll(pipe)
	if err != nil {
		log.Fatal("Cannot read output", loggertags.TagError, err)
	}

	return strings.Trim(string(bytes), " \t\n")
}

func goVersion() string {
	cmd := exec.Command("go", "version")
	return readOutput(cmd)
}

func gitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	return readOutput(cmd)
}

func gitCommit() string {
	cmd := exec.Command("git", "log", "-1", "--format=%h")
	return readOutput(cmd)
}

func gitStatus() string {
	cmd := exec.Command("git", "status", "--short")
	return strings.ReplaceAll(readOutput(cmd), "\n", "\\n")
}

func pegVersion() string {
	cmd := exec.Command("go", "tool", "peg", "-version")
	return strings.ReplaceAll(readOutput(cmd), "\n", "\\n")
}

func generate(packageName string) string {
	code := fmt.Sprintf(`const (
  GoVersion="%s"
  GitBranch="%s"
  GitCommit="%s"
  GitStatus="%s"
  PegVersion="%s"
  Version="Go: "+GoVersion+"; Git: "+GitBranch+"@"+GitCommit+"; Status:\n"+GitStatus+"\nPeg Version:"+PegVersion
)`,
		goVersion(),
		gitBranch(),
		gitCommit(),
		gitStatus(),
		pegVersion())

	return fmt.Sprintf(`package %s

// AUTO-GENERATED - DO NOT EDIT

%s`, packageName, code)
}

const outputFile = "build_info.go"

func main() {
	packageName := os.Args[1]
	log.Info("Generating build info")

	code := generate(packageName)
	err := os.WriteFile(outputFile, []byte(code), 0o777)
	if err != nil {
		log.Fatal("Cannot write to output folder", loggertags.TagError, err)
	}

	log.Info("Generated build info", loggertags.TagFileName, outputFile, "package", packageName)
}
