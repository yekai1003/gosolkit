package main

import (
	"os"
	"text/template"
)

type Inventory struct {
	Material string
	Count    uint
}

//1. 结构体与模版对应
//2. 结构体的字段必须与模版字段对应,大小写都一样
//3. 结构体字段顺序对模版没有影响

func main() {

	sweaters := Inventory{"wool", 17}
	tmpl, err := template.New("test").Parse("{{.Count}} of {{.Material}}\n")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, sweaters)
	if err != nil {
		panic(err)
	}
}
