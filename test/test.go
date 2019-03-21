package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("hello world")
	//abigen -sol pdbank.sol -pkg main -out pdbank.go
	cmd := exec.Command("abigen", "-sol", "../sol/pdbank.sol", "-pkg", "main", "-out", "pdbank.go")
	cmd.Run()
}
