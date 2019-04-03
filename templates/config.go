package templates

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
	SolidityPath   string //智能合约原路径
	GoPath         string //输出go代码路径
	AbiSH          string //处理abi的shell脚本路径
	BuildPath      string
	CodeName       string
	MainCodeName   string
	ContractName   string
	ConfigCodeName string
	ConfigTomlName string
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
