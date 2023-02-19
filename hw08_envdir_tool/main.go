package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	envDir, err := ReadDir(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	RunCmd(args[2:], envDir)
}
