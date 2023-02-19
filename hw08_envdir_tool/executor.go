package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	setEnvironment(env)
	returnCode = execCmd(cmd[0], cmd[1:])
	return
}

func execCmd(cmd string, args []string) int {
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	if err := command.Run(); err != nil {
		exitErr, ok := err.(*exec.ExitError)
		if ok {
			log.Println(err)
			return exitErr.ExitCode()
		} else {
			log.Fatal(err)
		}
	}
	return 0
}

func setEnvironment(env Environment) {
	for key, value := range env {
		if _, ok := os.LookupEnv(key); ok {
			if value.NeedRemove {
				err := os.Unsetenv(key)
				if err != nil {
					log.Println(err)
				}
				continue
			}
		}

		if !value.NeedRemove {
			err := os.Setenv(key, value.Value)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
