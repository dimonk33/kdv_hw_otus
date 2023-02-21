package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args
	envDir, err := ReadDir(args[1])
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(RunCmd(args[2:], envDir))
}
