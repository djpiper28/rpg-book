package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/log"
	loggertags "github.com/djpiper28/rpg-book/common/logger_tags"
)

func main() {
	args := []string{"--ts_out", "./src/lib/grpcClient/pb/", "--proto_path=../protos/"}

	// Expand the glob pattern for proto files
	protoFiles, err := filepath.Glob("../protos/*.proto")
	if err != nil {
		log.Error(os.Stderr, "Error expanding proto file glob", loggertags.TagError, err)
		os.Exit(1)
	}
	if len(protoFiles) == 0 {
		log.Error(os.Stderr, "No proto files found in ../protos/")
		os.Exit(1)
	}

	plugin := "./node_modules/.bin/protoc-gen-ts"
	if runtime.GOOS == "windows" {
		plugin += ".cmd"
	}

	cmdArgs := []string{"--plugin=protoc-gen-ts=" + plugin}
	cmdArgs = append(cmdArgs, args...)
	cmdArgs = append(cmdArgs, protoFiles...)

	cmd := exec.Command("protoc", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Info("Running command protoc", "args", cmdArgs)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "protoc command failed: %v\n", err)
		os.Exit(1)
	}
}
