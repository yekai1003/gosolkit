# Go语言打造以太坊智能合约测试框架

## 前言

### 这是什么？
这是一个基于go语言编写的，自动化测试以太坊智能合约的开发框架，使用此框架，可以自动化的部署合约，自动测试合约内的功能函数。你也发现了，本框架模拟的是truffle框架，但是完全是基于go语言编写，而且以太坊的客户端需要使用私链或者测试链。

### 可以学到什么？
- go语言访问以太坊
- 使用solc编译器编译合约
- 自动化部署合约
- abi文件解读
- go语言与智能合约调用
- golang模版编程
- golang处理toml配置文件
- shell脚本awk语言的运用
- go与命令行调用

本课程是指导开发者如何通过go语言来实现这样一个测试框架。

## 动起手来

### 环境准备 

- go语言开发环境
- 以太坊客户端-geth
- solc编译器

go语言环境安装本文不再详细说明，以下我们介绍其他环境安装。
#### geth安装
这个其实教程很多，搜索引擎都能搜到，也可以去官网查看！[官网安装说明](https://www.ethereum.org/cli)

### 安装solidity编译器

solidity智能合约需要编译，如果使用remix环境，在线编译器就帮我们做了，现在我们需要在命令行实现，那需要自己安装solc编译器，之后借助安装geth自带的abigen可执行程序，可以轻松的搞定solidity智能合约的编译，将sol文件编译为go文件。

下载solidity，这个由以太坊官方提供
```
git clone https://github.com/ethereum/solidity
```
此代码由c++编写实现，使用cmake进行编译，如果没有cmake需要安装一个。

- for mac
```
brew install cmake
```
- for ubuntu

```
sudo apt-get install cmake 
```
开始solidity代码

```
cd solidity
mkdir build
cd build
cmake .. && make
```
再将可执行文件solc拷贝到系统$PATH的某个路径下

```
cp solc/solc /usr/local/bin/
```

### 编译智能合约

当把工具准备好之后，我们就可以来尝试编译智能合约了。那么首先，我们得有一个智能合约。下面的例子是一个模拟银行业务的智能合约。



```
pragma solidity^0.5.0;

contract pdbank {
    address public  owner;
    mapping(address=>uint256) public balances;//纪录每个账户余额
    uint256 public totalAmount;
    string public bankName;//银行名称
    //构造函数
    constructor(string memory _bankName) public  {
        owner = msg.sender;
        bankName = _bankName;
    }
    //充值
    function deposit() public payable {
        //do nothing
        totalAmount += msg.value;
        balances[msg.sender] += msg.value;
    }
    //提现
    function withdraw(uint256 _amount) public payable {
        if(balances[msg.sender] > _amount) {
            balances[msg.sender]  -=  _amount;  
            msg.sender.transfer(_amount);
            totalAmount -= _amount;
        }
    }
}
```

如何编译它呢？查看以下abigen的帮助！

```
yekaideMacBook-Pro:~ yk$ abigen -h
Usage of abigen:
  -abi string
    	Path to the Ethereum contract ABI json to bind, - for STDIN
  -bin string
    	Path to the Ethereum contract bytecode (generate deploy method)
  -exc string
    	Comma separated types to exclude from binding
  -lang string
    	Destination language for the bindings (go, java, objc) (default "go")
  -out string
    	Output file for the generated binding (default = stdout)
  -pkg string
    	Package name to generate the binding into
  -sol string
    	Path to the Ethereum contract Solidity source to build and bind
  -solc string
    	Solidity compiler to use if source builds are requested (default "solc")
  -type string
    	Struct name for the binding (default = package name)

```

- abi选项是用于拿到abi文件后转换为go，但此时已经编译过了才会得到abi
- pkg 对应的输出go文件包名，很有用
- type 是对应输出合约结构类型，可以使用合约名字作为type
- out 指定输出文件
- sol 指定智能合约源码，我们本次要用到它

我们来新建一个工程：gosolkit


```
cd $GOPATH/src/
mkdir gosolkit
cd gosolkit
```

创建一个sol目录，用于存放智能合约


```
mkdir sol
```

创建pdbank.sol文件，编辑智能合约文件，将上述合约代码填入其中。


```
yekaideMacBook-Pro:sol yk$ cat pdbank.sol 
pragma solidity^0.5.0;

contract pdbank {
    address public  owner;
    mapping(address=>uint256) public balances;//纪录每个账户余额
    uint256 public totalAmount;
    string public bankName;//银行名称
    //构造函数
    constructor(string memory _bankName) public  {
        owner = msg.sender;
        bankName = _bankName;
    }
    //充值
    function deposit() public payable {
        //do nothing
        totalAmount += msg.value;
        balances[msg.sender] += msg.value;
    }
    //提现
    function withdraw(uint256 _amount) public payable {
        if(balances[msg.sender] > _amount) {
            balances[msg.sender]  -=  _amount;  
            msg.sender.transfer(_amount);
            totalAmount -= _amount;
        }
    }
}

```

编译代码，注意使用sol与pkg参数，type参数这里不启用（type在编译abi文件时使用，如果是源文件不需使用type）
```
yekaideMacBook-Pro:sol yk$ abigen -sol pdbank.sol -pkg main -out pdbank.go
yekaideMacBook-Pro:sol yk$ ls -lrt
total 48
-rw-r--r--  1 yk  staff    766  3 20 09:31 pdbank.sol
-rw-------  1 yk  staff  16915  3 20 09:33 pdbank.go

```
可以去欣赏一下这个go代码了！

### 调用智能合约

这里说的调用当然就是用go语言来调用智能合约。

首先来看发布，我们要在go语言代码中找到带Deploy开头的函数，很开心，我们找到了DeployPdbank。这个就是部署合约的函数，找到它，我们就可以部署合约了，那么就先研究这个函数如何使用！


```
func DeployPdbank(auth *bind.TransactOpts, backend bind.ContractBackend, _bankName string) (common.Address, *types.Transaction, *Pdbank, error)
```
参数说明
- auth 代表函数执行的身份，其实最重要的就是签名，因为我们知道部署合约是需要花费以太币的。
- backend 代表合约运行的环境或者说节点
- _bankName 就是我们合约里的参数了

返回值
- 合约地址
- 交易hash
- 部署后的合约对象
- 错误信息

我们先来考虑身份设置，需要账户地址以及私钥，私钥=keystore+password。

需要用到的工具包和函数：
- ethclient 可以提供上述的backend "github.com/ethereum/go-ethereum/ethclient"
- bind 可以提供TransactOpts，"github.com/ethereum/go-ethereum/accounts/abi/bind"
- common 可以提供数据类型的转换 "github.com/ethereum/go-ethereum/common"

编写部署合约的代码

```
func CallDeployPdbank() error {
	//链接到以太坊节点
	cli, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("failed to dial geht", err)
		return err
	}
	//创建身份，需要私钥
	auth, err := bind.NewTransactor(strings.NewReader(keydata), "123")
	if err != nil {
		fmt.Println("faild to NewTransactor", err)
		return err
	}
	//部署合约
	addr, ts, pb, err := DeployPdbank(auth, cli, "yekai")
	if err != nil {
		fmt.Println("failed to deploy pdbank", err)
		return err
	}
	//测试查询银行名
	name, _ := pb.BankName(nil)
	fmt.Println("addr=", addr.Hex(), "name=", name, ts.Hash().Hex())
	return nil
}
```

可以测试，执行结果发现：

```
yekaideMacBook-Pro:sol yk$ go run main.go pdbank.go 
addr= 0x625CC53E4660660eD94D90D263D1D793d453DE33 name=  0xc36f59cd82317ab14d49c5490cdbe6654a7f30ddc5eb5d208b277ac01fc2b4b3
```
也就是说name并没有获取成功，这是因为以太坊属于异步开发，刚刚部署不会立即生效，所以立即获取会失败！


再编写一个查询银行名字的
```
func CallBankName() error {
	cli, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("failed to dial geht", err)
		return err
	}
	defer cli.Close()
	//将之前部署的合约地址拿过来
	pb, err := NewPdbank(common.HexToAddress("0x625cc53e4660660ed94d90d263d1d793d453de33"), cli)
	//由于查看银行名称不需要消耗gas，所以不用指定身份
	name, _ := pb.BankName(nil)
	fmt.Println("bank name is", name)
	return err
}
```

再编写一个存款的函数调用

```
func CallDeposit() error {
	//链接到以太坊节点
	cli, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("failed to dial geht", err)
		return err
	}
	defer cli.Close()
	//创建身份，需要私钥
	auth, err := bind.NewTransactor(strings.NewReader(keydata), "123")
	if err != nil {
		fmt.Println("faild to NewTransactor", err)
		return err
	}
	pb, err := NewPdbank(common.HexToAddress("0x625cc53e4660660ed94d90d263d1d793d453de33"), cli)

	//存款需要携带金钱
	auth.Value = big.NewInt(10000001)
	//存款需要亮明身份
	ts, err := pb.Deposit(auth)
	if err != nil {
		fmt.Println("failed to deposit", err)
		return err
	}
	fmt.Println(ts.Hash().Hex())
	return err
}

```



