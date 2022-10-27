package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for name, value := range env {
		var err error

		if value.NeedRemove {
			err = os.Unsetenv(name)
		} else {
			err = os.Setenv(name, value.Value)
		}

		if err != nil {
			return 1
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return 1
	}

	return 0
}
