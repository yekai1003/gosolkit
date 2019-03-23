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

type InputsOutPuts struct {
	Name string
	Type string
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
	//定义模版
	deploy_temp, err := template.New("deploy").Parse(test_deploy_temp)
	if err != nil {
		fmt.Println("failed to template ", err)
		return err
	}
	var deploy_data DeployContractParams
	deploy_data.DeployName = "DeployPdbank"

	//对abi进行遍历处理
	for _, v := range abiInfos {
		if v.Type == "constructor" {
			// 如果是构造函数-部署函数
			deploy_data.DeployParams = "(auth,testclient"
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

	return nil
}

func Run() {
	Impl_run_code()

}
