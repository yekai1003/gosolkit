package main

import (
	"fmt"
	"gosolkit/templates"
	"os"
)

func Usage() {
	fmt.Printf("%s 1  -- compiler code\n", os.Args[0])
	fmt.Printf("%s 2  -- build test code\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		Usage()
		os.Exit(0)
	}
	if os.Args[1] == "1" {
		CompilerRun()
	} else if os.Args[1] == "2" {
		//build test code
		templates.Run()
	} else {
		Usage()
		os.Exit(0)
	}

}
