package main

import (
	"fmt"
	"golang.org/x/example/stringutil"
	"os"
)

const TestStr = "Hello, OTUS!"

func main() {
	_, _ = fmt.Fprintln(os.Stdout, stringutil.Reverse(TestStr))
}
