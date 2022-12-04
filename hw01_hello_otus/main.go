package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

const TestStr = "Hello, OTUS!"

func main() {
	fmt.Print(stringutil.Reverse(TestStr))
}
