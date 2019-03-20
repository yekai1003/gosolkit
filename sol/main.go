package main

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

///Users/yk/src/gowork/src/github.com/ethereum/go-ethereum

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

func CallBankName() error {
	//链接到以太坊,cli就是backend
	cli, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("failed to dial ", err)
		return err
	}
	//使用之前部署得到的合约地址
	//NewPdbank 合约调用的入口
	pd, err := NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), cli)
	if err != nil {
		fmt.Println("failed to NewPdbank", err)
		return err
	}
	//调用合约函数
	bkname, err := pd.BankName(nil)
	fmt.Println(bkname, err)

	return err
}

func CallDeposit() error {
	//1. 链接到以太坊
	cli, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		fmt.Println("failed to dial ", err)
		return err
	}
	//2. 构造函数入口 - 合约对象
	pd, err := NewPdbank(common.HexToAddress("0xD55E88D9156355C584982Db2C96dD1C2c63788C2"), cli)
	if err != nil {
		fmt.Println("failed to NewPdbank", err)
		return err
	}
	//3. 设置签名
	auth, err := bind.NewTransactor(strings.NewReader(keydata), "123")
	if err != nil {
		fmt.Println("failed to NewTransactor auth", err)
		return err
	}
	//4. 函数调用
	auth.Value = big.NewInt(100000011)
	ts, err := pd.Deposit(auth)
	if err != nil {
		fmt.Println("failed to Deposit ", err)
		return err
	}
	fmt.Println(ts.Hash().Hex())
	return err
}

func main() {
	CallDeploy()
	//CallBankName()
	//CallDeposit()
}
