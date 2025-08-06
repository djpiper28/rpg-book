package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	runtime "runtime"
)

func main() {
	args := []string{"--ts_out", "./src/lib/grpcClient/pb/", "--proto_path=../protos/"}

	// Expand the glob pattern for proto files
	protoFiles, err := filepath.Glob("../protos/*.proto")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error expanding proto file glob: %v\n", err)
		os.Exit(1)
	}
	if len(protoFiles) == 0 {
		fmt.Fprintln(os.Stderr, "No proto files found in ../protos/")
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

	fmt.Printf("Running command: protoc %v\n", cmdArgs)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "protoc command failed: %v\n", err)
		os.Exit(1)
	}
}
