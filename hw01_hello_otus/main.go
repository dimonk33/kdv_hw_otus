package main

import "golang.org/x/example/stringutil"

const TestStr = "Hello, OTUS!"

func main() {
	println(stringutil.Reverse(TestStr))
}
