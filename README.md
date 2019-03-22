## 第二课 智能合约自动化编译
### 前期内容回顾
之前我们的介绍的是如何通过solc编译智能合约，并且调用智能合约，本节我们继续实践，将智能合约的代码自动化编译以及abi文件生成搞定。

我们需要掌握什么技能呢？
 - go语言调用命令行
 - toml配置文件处理
- awk工具的使用

### go调用命令行

go调用命令行，我们使用exec包

- LookPath 可以判断一个可执行程序是否存在
- Command 创建一个命令行
- cmd.Run() 运行命令行，也可以使用Start()模式，可以去接收管道信息来得到程序返回结果
- 如果是一个shell脚本，那么可以用/bin/bash来启动

### toml配置文件处理
TOML的全称是Tom's Obvious, Minimal Language，因为它的作者是GitHub联合创始人Tom Preston-Werner。TOML 的目标是成为一个极简的配置文件格式。TOML 被设计成可以无歧义地被映射为哈希表，从而被多种语言解析。

[toml学习教程](https://github.com/BurntSushi/toml)
在使用的时候，记得要安装toml第三方包。

```
go get -u github.com/BurntSushi/toml
```
之后可以根据我们的需要，来编写配置文件，配置文件的目的仍然是为了让程序运行更灵活，而不应该成为我们的负担！

### awk工具使用

awk其实是一个语言，unix平台上处理文本的一种语言，其名称得自于它的创始人 Alfred Aho 、Peter Weinberger 和 Brian Kernighan 姓氏的首个字母。该语言的能力十分强大，可以支持字符串处理，打印等操作，工程中对于文本处理要求比较高的环节多会使用awk进行操作。

awk功能举例：


1. factory.txt 是一个工厂内目前产品库存情况，如果数量低于75，需要重新下订单，如何处理？
```
yekaideMacBook-Pro:awk yk$ cat factory.txt 
ProdA 70
ProdB 85
ProdC 74
```

示例如下：
```
awk '{if ($2 < 75) printf("%s reorder\n",$0);if ($2 >= 75) print $0}' factory.txt 
```

2. 查看本系统中shell是bash的用户名，并打印

```
cat /etc/passwd |grep bash|awk -F ":" '{print $1}'
```

3. awk处理合约的go文件，将abi信息截取处理存储到文本当中


```
awk '/const.+ABI = .+/{print substr($4,2,length($4)-2) }' pdbank.go > pdbank.abi
```

### 编写自动编译功能


main.go

```
package main

import (
	"fmt"
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
	} else if os.Args[1] == "1" {
		//build test code
	} else {
		Usage()
		os.Exit(0)
	}

}
```

接下来开始做填空题，也就是如何编译，我们先来实现。先编写扫描目录的代码，获取指定目录的sol文件，然后自动化的形成编译命令，送到命令行执行。

扫描指定目录的sol文件

```
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
```

编译函数
```
//编译一个智能合约
func CompilerOnece(solPath, solName, targetPath string) error {

	//xxx.sol - > xxx.go
	goName := strings.Replace(solName, ".sol", ".go", -1)

	cmd := exec.Command("abigen", "-sol", solPath+"/"+solName, "-pkg", targetPath, "-out", targetPath+"/"+goName)
	return cmd.Run()
}
```


构建abi函数，我们需要先用awk实现一个shell脚本，用来处理go文件的abi信息。
```
func BuildAbi(goCodeName string) error {
	abiName := strings.Replace(goCodeName, ".go", ".abi", -1)
	cmd := exec.Command("/bin/bash", "abi.sh", goCodeName, abiName)
	err := cmd.Run()
	fmt.Println("run BuildAbi ok!!", err)
	return nil
}
```
abi.sh

```
filename=$1
targetfile=$2
awk '/const.+ABI = .+/{print substr($4,2,length($4)-2) }' $filename > $targetfile

```

统一调用处理

```
func ParseRun() {
	solfiles, err := ParseDir("sol")
	fmt.Println(solfiles, err)
	for _, solfile := range solfiles {
		fmt.Println(solfile)
		codeName, err := Compiler(solfile, "sol", "contracts")
		if err != nil {
			fmt.Println("failed to complie code", err)
			return
		}
		err = BuildAbi(codeName)
		if err != nil {
			fmt.Println("failed to build abi", err)
			return
		}
	}
}
```
这样我们的基础工作完成了，但是代码不够完美，我们需要将部分写死的变量用配置文件来设置，所以再加入toml处理配置文件的部分。


添加config.tomls
```
[version]
auth = "yekai"
company = "pdj"
buildday = "2019-01-01"
ver = "1.0.0" # 版本

[common]
solidityPath = "sol"
goPath = "contracts"
abiSH = "./abi.sh"
```
添加config.go

```
package main

import (
	"log"

	"github.com/BurntSushi/toml"
)

type ServerConfig struct {
	Version VersionInfo
	Common  CommonInfo
}

type VersionInfo struct {
	Auth     string
	Company  string
	BuildDay string
	Ver      string
}

type CommonInfo struct {
	SolidityPath string //智能合约原路径
	GoPath       string //输出go代码路径
	AbiSH        string //处理abi的shell脚本路径
}

var ServConf ServerConfig

func init() {
	getConfig()
}

func getConfig() {
	var servConf ServerConfig
	_, err := toml.DecodeFile("config.toml", &servConf)
	if err != nil {
		log.Panic("faild to decodefile ", err)
	}
	ServConf = servConf
	//fmt.Println(servConf)
}


```
接下来替换原来的代码部分

```
package main

import (
	"fmt"
	"io/ioutil"
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

//构造abi
func BuildAbi(codePath string) error {
	//"contracts/pdbank.go"
	//xxx.sol - > xxx.abi
	abiName := strings.Replace(codePath, ".sol", ".abi", -1)
	goName := strings.Replace(codePath, ".sol", ".go", -1)

	cmd := exec.Command(ServConf.Common.AbiSH, goName, abiName)
	return cmd.Run()
}

//扫描目录，获得全部的文件
func CompilerRun() error {
	infos, err := ioutil.ReadDir(ServConf.Common.SolidityPath)
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
			err = CompilerOnece(ServConf.Common.SolidityPath, v.Name(), ServConf.Common.GoPath)
			if err != nil {
				fmt.Println("call ompilerOnece err", err)
				break
			}
			//创建abi
			err = BuildAbi(ServConf.Common.GoPath + "/" + v.Name())
			if err != nil {
				fmt.Println("call BuildAbi err", err)
				break
			}
		}
	}
	return err
}

```

### 之前内容
[第一步：solidity合约编译与调用](https://github.com/yekai1003/gosolkit/blob/level1/README.md)

