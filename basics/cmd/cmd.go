package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func absPath(cmd string) (string, error) {
	cmdAbsPath, err := exec.LookPath(cmd)
	if err != nil {
		return "", err
	}

	return cmdAbsPath, nil
}

func buildCommand(cmd string, args ...string) (*exec.Cmd, error) {
	cmdAbsPath, err := absPath(cmd)
	if err != nil {
		return nil, err
	}

	command := exec.Command(cmdAbsPath)
	command.Args = append(command.Args, args...)
	return command, nil
}

func RunCommand(cmd string, args ...string) ([]string, error) {
	command, err := buildCommand(cmd, args...)
	if err != nil {
		return nil, err
	}

	output, err := command.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(string(output)), "\n"), nil
}

func main() {
	output, err := RunCommand("ls", "/tmp")
	if err != nil {
		fmt.Printf("Error %v\n", err)
	} else {
		fmt.Printf("%s", strings.Join(output, "\n"))
	}
}
