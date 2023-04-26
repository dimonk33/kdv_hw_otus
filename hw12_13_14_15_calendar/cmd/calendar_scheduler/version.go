package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	release   = "1.0"
	buildDate = "2020-04-26"
	gitHash   = "UNKNOWN"
)

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}
