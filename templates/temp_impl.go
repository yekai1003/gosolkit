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

/*
{

	"type": "function"
}
*/

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
	//3. 写入部署合约代码
	//定义模版
	deploy_temp, err := template.New("deploy").Parse(test_deploy_temp)
	if err != nil {
		fmt.Println("failed to template ", err)
		return err
	}
	var deploy_data DeployContractParams
	deploy_data.DeployName = "DeployPdbank"
	deploy_data.DeployParams = "(auth,testclient,\"yekai\")" //自动生成
	//模版的执行
	err = deploy_temp.Execute(outfile, &deploy_data)
	if err != nil {
		fmt.Println("failed to template Execute ", err)
		return err
	}

	return nil
}

func Run() {
	Impl_run_code()
	infos, _ := readAbi("contracts/pdbank.abi")
	fmt.Println(infos)
}
