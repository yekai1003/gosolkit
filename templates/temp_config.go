package templates

const config_build_temp = `
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
	ConnStr      string //以太坊节点地址
	KeyStoreDir  string //key存放路径
	ContractAddr string //合约地址
	DeployAddr   string
	DeployPass   string
	TestAddr     string
	TestPass     string
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
`

const config_toml_temp = `
[version]
auth = "{{.Auth}}"
company = "{{.Company}}"
buildday = "{{.BuildDay}}"
ver = "{{.Ver}}" # 版本

[common]
connStr = "http://localhost:8545"
keyStoreDir = "/Users/yk/ethdev/data/keystore"
contractAddr = "0xD55E88D9156355C584982Db2C96dD1C2c63788C2" #部署之后得到
deployAddr = "0x791443d21a76e16cc510b6b1684344d2a5ce751c"
deployPass = "123"
testAddr = "0x791443d21a76e16cc510b6b1684344d2a5ce751c"
testPass = "123"
`
