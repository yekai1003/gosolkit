package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

//编译一个智能合约
func CompilerOnece(solPath, solName, targetPath string) error {

	//xxx.sol - > xxx.go
	goName := strings.Replace(solName, ".sol", ".go", -1)

	cmd := exec.Command("abigen", "-sol", solPath+"/"+solName, "-pkg", targetPath, "-out", targetPath+"/"+goName)
	return cmd.Run()
}

//扫描目录，获得全部的文件
func CompilerRun() error {
	infos, err := ioutil.ReadDir("sol")
	if err != nil {
		fmt.Println("failed to readdir ", err)
		return err
	}
	for _, v := range infos {

		//后4位位.sol
		strNameRune := []rune(v.Name())
		strfix := string(strNameRune[len(strNameRune)-4:])
		if strfix == ".sol" && !v.IsDir() {
			fmt.Println(v.IsDir(), v.Name(), v.Size(), "ok")
			err = CompilerOnece("sol", v.Name(), "contracts")
			if err != nil {
				fmt.Println("call ompilerOnece err", err)
				break
			}
		}
	}
	return err
}
