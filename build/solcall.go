
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

func CallTotalAmount() error {

	//使用之前部署得到的合约地址
	instance, err := contracts.NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testClient)
	if err != nil {
		fmt.Println("failed to instance contract", err)
		return err
	}
	//调用合约函数
	data0,err := instance.TotalAmount(nil)
	fmt.Println(data0,err)

	return err
}

func CallBankName() error {

	//使用之前部署得到的合约地址
	instance, err := contracts.NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testClient)
	if err != nil {
		fmt.Println("failed to instance contract", err)
		return err
	}
	//调用合约函数
	data0,err := instance.BankName(nil)
	fmt.Println(data0,err)

	return err
}

func CallBalances() error {

	//使用之前部署得到的合约地址
	instance, err := contracts.NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testClient)
	if err != nil {
		fmt.Println("failed to instance contract", err)
		return err
	}
	//调用合约函数
	data0,err := instance.Balances(nil,common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"))
	fmt.Println(data0,err)

	return err
}

func CallWithdraw() error {

	//2. 构造函数入口 - 合约对象
	instance, err := contracts.NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testClient)
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
	ts, err := instance.Withdraw(auth,big.NewInt(1000))
	if err != nil {
		fmt.Println("failed to Deposit ", err)
		return err
	}
	fmt.Println(ts.Hash().Hex())
	return err
}

func CallOwner() error {

	//使用之前部署得到的合约地址
	instance, err := contracts.NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testClient)
	if err != nil {
		fmt.Println("failed to instance contract", err)
		return err
	}
	//调用合约函数
	data0,err := instance.Owner(nil)
	fmt.Println(data0,err)

	return err
}

func CallDeposit() error {

	//2. 构造函数入口 - 合约对象
	instance, err := contracts.NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), testClient)
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
	ts, err := instance.Deposit(auth)
	if err != nil {
		fmt.Println("failed to Deposit ", err)
		return err
	}
	fmt.Println(ts.Hash().Hex())
	return err
}

func CallDeploy() error {
	//创建身份,需要私钥= pass+keystore文件
	auth, err := MakeAuth("0x791443d21a76e16cc510b6b1684344d2a5ce751c", "123")
	if err != nil {
		fmt.Println("failed to MakeAuth auth", err)
		return err
	}
	addr, ts, _, err := contracts.DeployPdbank(auth,testClient,"yekai")
	if err != nil {
		fmt.Println("failed to DeployPdbank", err)
		return err
	}

	fmt.Println("addr=", addr.Hex(), "hash=", ts.Hash().Hex())
	return err
}
