package main

import (
	"os"
	"os/exec"
)

const (
	ExitCodeSuccess = 0
	ExitCodeFail    = 1
)

func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return ExitCodeFail
	}

	cmdName, args := cmd[0], cmd[1:]

	proc := exec.Command(cmdName, args...)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin

	for key, value := range env {
		_, exists := os.LookupEnv(key)

		if exists {
			if err := os.Unsetenv(key); err != nil {
				return ExitCodeFail
			}
		}

		if value.NeedRemove {
			continue
		}

		if err := os.Setenv(key, value.Value); err != nil {
			return ExitCodeFail
		}
	}

	if err := proc.Run(); err != nil {
		return ExitCodeFail
	}

	return ExitCodeSuccess
}
