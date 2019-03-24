package templates

const test_run_main_temp = `
package main

import (
	"fmt"
	"os"
)

var funcNames = []string{%s}
`

//1. 提供一个命令行帮助
const test_build_main_temp = `
func Usage() {
	fmt.Printf("%s 1 -- deploy\n", os.Args[0])
	num := 2
	for _, v := range funcNames {
		fmt.Printf("%s %d -- %s\n", os.Args[0], num, v)
		num++
	}
}


func main() {
	if len(os.Args) < 2 {
		Usage()
		os.Exit(0)
	}
	if os.Args[1] == "1" {
		CallDeploy()
	}{{range.}} else if os.Args[1] == "{{.Num}}" {
		Call{{.FuncName}}()
	} {{end}}
}

`
