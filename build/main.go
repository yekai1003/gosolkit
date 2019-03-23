package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello world")
	if os.Args[1] == "2" {
		CallDeploy()
	}
}
