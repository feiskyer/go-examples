package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/docker/docker/pkg/stdcopy"
)

func buildCommand(args ...string) *exec.Cmd {
	hyperBinAbsPath, err := exec.LookPath("hyper")
	if err != nil {
		return nil
	}

	cmd := exec.Command(hyperBinAbsPath)
	cmd.Args = append(cmd.Args, args...)
	return cmd
}

func main() {
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr

	args := append([]string{}, "attach", "1fe72452ff9ab254ba90a7ec899f84b54bfe0d8198563810d24adb6b792e65f7")
	command := buildCommand(args...)

	if stdin != nil {
		r, w, err := os.Pipe()
		if err != nil {
			fmt.Printf("Pipe: %v", err)
		}
		go io.Copy(w, stdin)

		command.Stdin = r
	}
	if stdout != nil {
		command.Stdout = stdout
	}
	if stderr != nil {
		command.Stderr = stderr
	}

	err := command.Run()
	fmt.Printf("Result: %v", err)
}
