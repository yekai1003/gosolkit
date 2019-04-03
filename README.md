## 第三课 智能合约自动化测试
### 之前课程回顾
我们之前介绍了go语言调用exec处理命令行，介绍了toml配置文件的处理，以及awk处理文本文件获得ABI信息。我们的代码算是完成了从智能合约到go语言的自动编译，同时也可以自动提取到ABI信息。

具体可以参考：

[第一课 go语言与智能合约调用的来龙去脉](https://github.com/yekai1003/gosolkit/blob/level1/README.md)

[第二课 智能合约自动化编译](https://github.com/yekai1003/gosolkit/blob/level2/README.md)


###  本节主要工作
- go语言模版编程
- 目标代码生成
- 目标代码的调用代码生成
- 配置文件的自动化生成

### 原有代码优化

想要自动化生成测试代码，首先你要知道目标代码长什么样儿，也就是之前我们写的部署合约，调用合约的代码。仔细思考，其实会发现分了两层内容，第一层内容就是针对每一个合约函数都要有一个调用函数生成，假设合约函数为X，我们把调用该函数的函数叫CallX，我们想要完成测试，又需要在main函数内能找到调用CallX的入口。先一步一步来，先完成一个小目标，那就是生成这些个CallX，
再来梳理一下我们之前调用合约的代码。


```
const keydata = `{"address":"791443d21a76e16cc510b6b1684344d2a5ce751c","crypto":{"cipher":"aes-128-ctr","ciphertext":"bbccbf9deb8c907d9f245767fffb57880c4cfd265dde9372d7278a8e963043bd","cipherparams":{"iv":"95e8f925fe0f3460f0ca3ccebc481b14"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"2036300fae07d954a10e70a2f87876fa198b310dc16a72f3fab3265978e7d798"},"mac":"009ce0ecd11d5f563d2f7bd41eb957f1ca3b517f31653dd18fac20e77b7feb5c"},"id":"7f42641c-580a-4398-b00c-74ef2710bcae","version":3}`

func CallDeploy() error {
	//链接到以太坊,cli就是backend
	cli, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("failed to dial ", err)
		return err
	}
	//创建身份,需要私钥= pass+keystore文件
	auth, err := bind.NewTransactor(strings.NewReader(keydata), "123")
	if err != nil {
		fmt.Println("failed to NewTransactor auth", err)
		return err
	}
	auth.GasLimit = 400000
	addr, ts, pd, err := DeployPdbank(auth, cli, "yekai233")
	if err != nil {
		fmt.Println("failed to DeployPdbank", err)
		return err
	}
	bkname, _ := pd.BankName(nil)
	fmt.Println("addr=", addr.Hex(), "bkname=", bkname, "hash=", ts.Hash().Hex())
	return err
}
```

对于keydata来说，我们不能写死了，需要从keystore文件读取这部分数据。http://localhost:8545这样的连接地址最好通过配置文件来搞定，但思路要清楚的是，我们这个代码是将来要执行的，这个代码它需要一个自己的配置文件，并非是我们当前工程的配置文件，直白一点的说就是我们需要给未来的代码宝宝们自动的生成一个配置文件。记得有这么回事儿，我们先不管这一点，先以能完成任务为前提。


我们还是来单独形成一个函数，用于专门形成签名，它需要读取keystore文件，这个文件可以根据对应的账户地址到以太坊keystore目录得到，这样就容易多了！
```
//设置签名
func MakeAuth(addr, pass string) (*bind.TransactOpts, error) {
	keystorePath  :=  "{{.Keydir}}"
	fileName, err := GetFileName(string([]rune(addr)[2:]), keystorePath)
	if err != nil {
		fmt.Println("failed to GetFileName", err)
		return nil, err
	}

	file, err := os.Open(keystorePath + "/" + fileName)
	if err != nil {
		fmt.Println("failed to open file ", err)
		return nil, err
	}
	auth, err := bind.NewTransactor(file, pass)
	if err != nil {
		fmt.Println("failed to NewTransactor  ", err)
		return nil, err
	}
	return auth, err
}
func GetFileName(address, dirname string) (string, error) {

	data, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println("read dir err", err)
		return "", err
	}
	for _, v := range data {
		if strings.Index(v.Name(), address) > 0 {
			//代表找到文件
			return v.Name(), nil
		}
	}

	return "", nil
}
```
于是，我们将原来写的测试代码进行完善一下！

此外，连接这里也规范一下，使用init函数来搞定！  

```
var testclient *ethclient.Client

func init() {
	cli, err := CreateCli("http://localhost:8545")
	if err != nil {
		log.Panic("failed to connect to eth", err)
	}
	testclient = cli
}
```


```
func CallBankName() (error) {
	instance, err := NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testclient)
	if err != nil {
		fmt.Println("failed to get contract instance", err)
		return err
	}
	data,err := instance.BankName(nil)
	if err != nil {
		fmt.Println("failed to get Balances", err)
		return err
	}
	fmt.Println(data,err)
	return nil
}


func CallWithdraw(addr, pass string) (*types.Transaction, error) {

	instance, err := NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testclient)
	if err != nil {
		fmt.Println("failed to get contract instance", err)
		return nil, err
	}
	auth, err := MakeAuth(addr, pass)
	if err != nil {
		fmt.Println("failed to makeAuth", err)
		return nil, err
	}
	auth.Value = big.NewInt(0)
	ts,err := instance.Withdraw(auth,big.NewInt(10000))
	if err != nil {
		fmt.Println("failed to call ", err)
		return nil, err
	}
	fmt.Println(ts.ChainId(), ts.Hash().Hex(), ts.Nonce())
	return ts , err
}
```

这样代码看上去舒服多了，接下来考虑目标代码的生成办法。

### go语言与模版编程

go语言模版编程需要用到template包，go语言的模版呢实际上也提供了两大类模版处理。一类是文本的，一类是html的。鉴于我们的目标代码是go语言，所以本次使用基于文本的模版处理。


下面是一个最简单的模版处理的例子，{{.Count}}和{{.Material}}相当于是这个模版要填的两个空。
```
type Inventory struct {
	Material string
	Count    uint
}
sweaters := Inventory{"wool", 17}
tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
if err != nil { panic(err) }
err = tmpl.Execute(os.Stdout, sweaters)
if err != nil { panic(err) }
```
- New 方法创建一个模版
- Parse  是指定模版的内容
- Execute 是执行，两个参数分别是写入的方向以及模版实例化的内容

一个复杂一些的例子
```
// Define a template.
const letter = `
Dear {{.Name}},
{{if .Attended}}
It was a pleasure to see you at the wedding.{{else}}
It is a shame you couldn't make it to the wedding.{{end}}
{{with .Gift}}Thank you for the lovely {{.}}.
{{end}}
Best wishes,
Josie
`

// Prepare some data to insert into the template.
type Recipient struct {
    Name, Gift string
    Attended   bool
}
var recipients = []Recipient{
    {"Aunt Mildred", "bone china tea set", true},
    {"Uncle John", "moleskin pants", false},
    {"Cousin Rodney", "", false},
}

// Create a new template and parse the letter into it.
t := template.Must(template.New("letter").Parse(letter))

// Execute the template for each recipient.
for _, r := range recipients {
    err := t.Execute(os.Stdout, r)
    if err != nil {
        log.Println("executing template:", err)
    }
}
```

我们可以发现一些特点：
-  {{.Name}}这样的pipeline（官方学名）一定能在结构体中找到这个名字
-  可以使用if和with做流程处理 

于是把我们的目标代码，可以制作成定制化的模版！

创建一个用于存放模版定义的文件,template_code.go


```
package templates

const Main_tmpl = `package main

import (
	"fmt"
	"log"
	"os"

	"gosol/contracts"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

var testclient *ethclient.Client

func init() {
	cli, err := CreateCli("{{.Connstr}}")
	if err != nil {
		log.Panic("failed to connect to eth", err)
	}
	testclient = cli
}

func GetFileName(address, dirname string) (string, error) {

	data, err := ioutil.ReadDir(dirname)
	if err != nil {
		fmt.Println("read dir err", err)
		return "", err
	}
	for _, v := range data {
		if strings.Index(v.Name(), address) > 0 {
			//代表找到文件
			return v.Name(), nil
		}
	}

	return "", nil
}

//创建链接
func CreateCli(connstr string) (*ethclient.Client, error) {
	cli, err := ethclient.Dial(connstr)
	if err != nil {
		fmt.Println("failed to dial provide", err)
		return nil, err
	}
	return cli, err
}

//设置签名
func MakeAuth(addr, pass string) (*bind.TransactOpts, error) {
	keystorePath  :=  "{{.Keydir}}"
	fileName, err := GetFileName(string([]rune(addr)[2:]), keystorePath)
	if err != nil {
		fmt.Println("failed to GetFileName", err)
		return nil, err
	}

	file, err := os.Open(keystorePath + "/" + fileName)
	if err != nil {
		fmt.Println("failed to open file ", err)
		return nil, err
	}
	auth, err := bind.NewTransactor(file, pass)
	if err != nil {
		fmt.Println("failed to NewTransactor  ", err)
		return nil, err
	}
	return auth, err
}
`

const Deploy_sol_tmpl = `
func Deploy{{.ContractName}}() (common.Address, error) {
	auth, err := MakeAuth("{{.FromAddr}}", "{{.Pass}}")
	if err != nil {
		fmt.Println("failed to makeAuth", err)
		return common.HexToAddress(""), err
	}

	//common.Address, *types.Transaction, *Pdbank, error
	contractaddr, ts, _, err := contracts.{{.CallFunc}}
	if err != nil {
		fmt.Println("failed to deloy ",err)
		return common.HexToAddress(""), err
	}
	fmt.Println(ts.ChainId(), ts.Hash().Hex(), ts.Nonce())
	fmt.Println(contractaddr.Hex())
	return contractaddr, err
}

`
```
先搞定部署合约代码自动生成部分！


```
//1. 写到哪
	outfile, err := os.OpenFile("build/solcall.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}
	defer outfile.Close()
	//2. 写什么
	_, err = outfile.WriteString(test_main_temp)
	if err != nil {
		fmt.Println("failed to write ", err)
		return err
	}
	// 读取abi文件信息
	abiInfos, err := readAbi("contracts/pdbank.abi")
	if err != nil {
		fmt.Println("failed to read abi", err)
		return err
	}
	//fmt.Println(infos)

	//3. 写入部署合约代码
	//定义部署模版
	deploy_temp, err := template.New("deploy").Parse(test_deploy_temp)
	if err != nil {
		fmt.Println("failed to template deploy", err)
		return err
	}
	var deploy_data DeployContractParams
	deploy_data.DeployName = "DeployPdbank"
	
	for _, v := range abiInfos {
		v.Name = strings.Title(v.Name) //标题优化，首字母大写, hello world - > Hello World
		if v.Type == "constructor" {
			// 如果是构造函数-部署函数
			deploy_data.DeployParams = "(auth,testClient"
			for _, vv := range v.Inputs {
				//需要根据输入数据类型来判断如何处理:string,address,uint256
				if vv.Type == "address" {
					deploy_data.DeployParams += ",common.HexToAddress(\"0xD55E88D9156355C584982Db2C96dD1C2c63788C2\")"
				} else if vv.Type == "uint256" {
					deploy_data.DeployParams += ",big.NewInt(1000)"
				} else if vv.Type == "string" {
					deploy_data.DeployParams += ",\"yekai\""
				}

			}
			deploy_data.DeployParams += ")"
			//模版的执行
			err = deploy_temp.Execute(outfile, &deploy_data)
			if err != nil {
				fmt.Println("failed to template Execute ", err)
				return err
			}
		}
	}
```

接下来的代码还是相同的套路。
完整代码如下；



```
//temp_org.go
package templates

const test_main_temp = `
package main

import (
	"fmt"
	"gosolkit/contracts"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var testClient *ethclient.Client

func init() {
	cli, err := Connect("http://localhost:8545")
	if err != nil {
		log.Fatalln("failed to connect to eth", err)
	}
	testClient = cli
}

func Connect(connstr string) (*ethclient.Client, error) {
	return ethclient.Dial(connstr)
}

//签名函数
func MakeAuth(addr, pass string) (*bind.TransactOpts, error) {
	//1. 根据addr找到keystore目录下的文件
	keyDir := "/Users/yk/ethdev/data/keystore"
	infos, err := ioutil.ReadDir(keyDir)
	if err != nil {
		fmt.Println("failed to readdir", err)
		return nil, err
	}
	//UTC--2019-03-16T13-00-48.032030904Z--791443d21a76e16cc510b6b1684344d2a5ce751c
	//0x791443d21a76e16cc510b6b1684344d2a5ce751c
	strAddr := ([]rune(addr))[2:]
	for _, v := range infos {
		strVname := []rune(v.Name())
		if len(strVname) > len(strAddr) {
			strVname2 := strVname[len(strVname)-len(strAddr):]
			if strings.EqualFold(string(strAddr), string(strVname2)) {
				//找到了匹配的文件
				//fmt.Println(addr, v.Name())
				//2. 做签名
				reader, err := os.Open(keyDir + "/" + v.Name())
				if err != nil {
					fmt.Println("failed to open file", err)
					return nil, err
				}
				defer reader.Close()
				auth, err := bind.NewTransactor(reader, pass)
				if err != nil {
					fmt.Println("failed to NewTransactor auth", err)
					return nil, err
				}
				return auth, err
			}
		}
	}

	return nil, nil
}
`
const test_deploy_temp = `
func CallDeploy() error {
	//创建身份,需要私钥= pass+keystore文件
	auth, err := MakeAuth("0x791443d21a76e16cc510b6b1684344d2a5ce751c", "123")
	if err != nil {
		fmt.Println("failed to MakeAuth auth", err)
		return err
	}
	addr, ts, _, err := contracts.{{.DeployName}}{{.DeployParams}}
	if err != nil {
		fmt.Println("failed to DeployPdbank", err)
		return err
	}

	fmt.Println("addr=", addr.Hex(), "hash=", ts.Hash().Hex())
	return err
}
`

const test_nogas_temp = `
func Call{{.FuncName}}() error {

	//使用之前部署得到的合约地址
	instance, err := contracts.{{.NewContractName}}(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testClient)
	if err != nil {
		fmt.Println("failed to instance contract", err)
		return err
	}
	//调用合约函数
	{{.OutParams}} := instance.{{.FuncName}}{{.InputParams}}
	fmt.Println({{.OutParams}})

	return err
}
`

const test_gas_temp = `
func Call{{.FuncName}}() error {

	//2. 构造函数入口 - 合约对象
	instance, err := contracts.{{.NewContractName}}(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testClient)
	if err != nil {
		fmt.Println("failed to contract instance", err)
		return err
	}
	//3. 设置签名
	auth, err := MakeAuth("0x791443d21a76e16cc510b6b1684344d2a5ce751c", "123")
	if err != nil {
		fmt.Println("failed to MakeAuth auth", err)
		return err
	}
	//4. 函数调用
	auth.Value = big.NewInt(0)
	ts, err := instance.{{.FuncName}}{{.InputParams}}
	if err != nil {
		fmt.Println("failed to Deposit ", err)
		return err
	}
	fmt.Println(ts.Hash().Hex())
	return err
}
`
```


```
//temp_main.go
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
```



```
//temp_impl.go
package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type DeployContractParams struct {
	DeployName   string
	DeployParams string
}

//无gas函数调用
type FuncNoGasParams struct {
	FuncName        string
	NewContractName string
	OutParams       string
	InputParams     string
}

//有gas函数调用
type FuncGasParams struct {
	FuncName        string
	NewContractName string
	InputParams     string
}

type InputsOutPuts struct {
	Name string
	Type string
}

type FuncInfo struct {
	FuncName string
	Num      int
}

type AbiInfo struct {
	Constant        bool
	Inputs          []InputsOutPuts
	Name            string
	Outputs         []InputsOutPuts
	Payable         bool
	StateMutability string
	Type            string
}

func readAbi(abifile string) ([]AbiInfo, error) {
	file, err := os.Open(abifile)
	if err != nil {
		fmt.Println("failed to open file ", err)
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("failed to read abi", err)
		return nil, err
	}
	var abiInfos []AbiInfo
	strdata := strings.Replace(string(data), "\\", "", -1)
	err = json.Unmarshal([]byte(strdata), &abiInfos)
	if err != nil {
		fmt.Println("failed to Unmarshal abi", err)
		return nil, err
	}
	return abiInfos, err
}

func Impl_run_code() error {
	//1. 写到哪
	outfile, err := os.OpenFile("build/solcall.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}
	defer outfile.Close()
	//2. 写什么
	_, err = outfile.WriteString(test_main_temp)
	if err != nil {
		fmt.Println("failed to write ", err)
		return err
	}
	// 读取abi文件信息
	abiInfos, err := readAbi("contracts/pdbank.abi")
	if err != nil {
		fmt.Println("failed to read abi", err)
		return err
	}
	//fmt.Println(infos)

	//3. 写入部署合约代码
	//定义部署模版
	deploy_temp, err := template.New("deploy").Parse(test_deploy_temp)
	if err != nil {
		fmt.Println("failed to template deploy", err)
		return err
	}
	var deploy_data DeployContractParams
	deploy_data.DeployName = "DeployPdbank"

	//定义nogas函数的模版
	nogas_temp, err := template.New("nogas").Parse(test_nogas_temp)
	if err != nil {
		fmt.Println("failed to template nogas_temp", err)
		return err
	}

	var func_nogas_data FuncNoGasParams
	func_nogas_data.NewContractName = "NewPdbank"

	//定义有gas模版
	hasgas_temp, err := template.New("hasgas").Parse(test_gas_temp)
	if err != nil {
		fmt.Println("failed to template hasgas_temp", err)
		return err
	}

	var func_gas_data FuncGasParams
	func_gas_data.NewContractName = "NewPdbank"

	//对abi进行遍历处理
	for _, v := range abiInfos {
		v.Name = strings.Title(v.Name) //标题优化，首字母大写, hello world - > Hello World
		if v.Type == "constructor" {
			// 如果是构造函数-部署函数
			deploy_data.DeployParams = "(auth,testClient"
			for _, vv := range v.Inputs {
				//需要根据输入数据类型来判断如何处理:string,address,uint256
				if vv.Type == "address" {
					deploy_data.DeployParams += ",common.HexToAddress(\"0xD55E88D9156355C584982Db2C96dD1C2c63788C2\")"
				} else if vv.Type == "uint256" {
					deploy_data.DeployParams += ",big.NewInt(1000)"
				} else if vv.Type == "string" {
					deploy_data.DeployParams += ",\"yekai\""
				}

			}
			deploy_data.DeployParams += ")"
			//模版的执行
			err = deploy_temp.Execute(outfile, &deploy_data)
			if err != nil {
				fmt.Println("failed to template Execute ", err)
				return err
			}
		} else {
			//处理其他函数
			if len(v.Outputs) > 0 {
				//不需要gas函数
				func_nogas_data.FuncName = v.Name

				func_nogas_data.InputParams = "(nil"
				for _, vv := range v.Inputs {
					//需要根据输入数据类型来判断如何处理:string,address,uint256
					if vv.Type == "address" {
						func_nogas_data.InputParams += ",common.HexToAddress(\"0xD55E88D9156355C584982Db2C96dD1C2c63788C2\")"
					} else if vv.Type == "uint256" {
						func_nogas_data.InputParams += ",big.NewInt(1000)"
					} else if vv.Type == "string" {
						func_nogas_data.InputParams += ",\"yekai\""
					}

				}
				func_nogas_data.InputParams += ")"
				//输入参数
				num := 0
				strOutPuts := ""
				for _, _ = range v.Outputs {
					strOutPuts = fmt.Sprintf("%sdata%d,", strOutPuts, num)
					num++
				}
				strOutPuts += "err"
				func_nogas_data.OutParams = strOutPuts

				//模版的执行
				err = nogas_temp.Execute(outfile, &func_nogas_data)
				if err != nil {
					fmt.Println("failed to template nogas Execute ", err)
					return err
				}
			} else {
				//需要消耗gas
				func_gas_data.FuncName = v.Name
				func_gas_data.InputParams = "(auth"
				for _, vv := range v.Inputs {
					//需要根据输入数据类型来判断如何处理:string,address,uint256
					if vv.Type == "address" {
						func_gas_data.InputParams += ",common.HexToAddress(\"0xD55E88D9156355C584982Db2C96dD1C2c63788C2\")"
					} else if vv.Type == "uint256" {
						func_gas_data.InputParams += ",big.NewInt(1000)"
					} else if vv.Type == "string" {
						func_gas_data.InputParams += ",\"yekai\""
					}

				}
				func_gas_data.InputParams += ")"
				//模版的执行
				err = hasgas_temp.Execute(outfile, &func_gas_data)
				if err != nil {
					fmt.Println("failed to template hasgas Execute ", err)
					return err
				}
			}
		}
	}

	return nil
}

func Impl_main_code() error {
	//1. 写到哪
	outfile, err := os.OpenFile("build/main.go", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}
	defer outfile.Close()
	// 读取abi文件信息
	abiInfos, err := readAbi("contracts/pdbank.abi")
	if err != nil {
		fmt.Println("failed to read abi", err)
		return err
	}
	funcNames := ""
	//"abc","def","
	num := 0
	var funcInfos []FuncInfo
	var funcInfo FuncInfo
	// 2- 第一个函数
	for _, v := range abiInfos {
		if v.Type != "constructor" {
			if num == 0 {
				//第一个
				funcNames += fmt.Sprintf(`"%s"`, v.Name)
			} else {
				funcNames += fmt.Sprintf(`,"%s"`, v.Name)
			}
			num++
			funcInfo.FuncName = strings.Title(v.Name)
			funcInfo.Num = num + 1
			funcInfos = append(funcInfos, funcInfo)
		}
	}
	main_str1 := fmt.Sprintf(test_run_main_temp, funcNames)
	_, err = outfile.WriteString(main_str1)
	if err != nil {
		fmt.Println("failed to write to main.go", err)
		return err
	}

	//建立一个模版，输出内容
	main_temp, err := template.New("main").Parse(test_build_main_temp)
	if err != nil {
		fmt.Println("failed to template main", err)
		return err
	}
	err = main_temp.Execute(outfile, funcInfos)
	if err != nil {
		fmt.Println("failed to Execute main", err)
		return err
	}
	return err
}

func Run() {
	Impl_run_code()
	Impl_main_code()
}
```

在main函数内增加此部分调用


```
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
```
