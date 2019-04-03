
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
