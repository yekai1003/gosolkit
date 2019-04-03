package templates

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"
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
	outfile, err := os.OpenFile(ServConf.Common.BuildPath+"/"+ServConf.Common.CodeName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
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
	abiInfos, err := readAbi(ServConf.Common.GoPath + "/" + ServConf.Common.ContractName + ".abi")
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
	func_nogas_data.NewContractName = "New" + strings.Title(ServConf.Common.ContractName)

	//定义有gas模版
	hasgas_temp, err := template.New("hasgas").Parse(test_gas_temp)
	if err != nil {
		fmt.Println("failed to template hasgas_temp", err)
		return err
	}

	var func_gas_data FuncGasParams
	func_gas_data.NewContractName = "New" + strings.Title(ServConf.Common.ContractName)

	//对abi进行遍历处理
	for _, v := range abiInfos {
		v.Name = strings.Title(v.Name) //标题优化，首字母大写, hello world - > Hello World
		if v.Type == "constructor" {
			// 如果是构造函数-部署函数
			deploy_data.DeployParams = "(auth,testClient"
			for _, vv := range v.Inputs {
				//需要根据输入数据类型来判断如何处理:string,address,uint256
				if vv.Type == "address" {
					deploy_data.DeployParams += ",common.HexToAddress(ServConf.Common.TestAddr)"
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
						func_nogas_data.InputParams += ",common.HexToAddress(ServConf.Common.TestAddr)"
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
						func_gas_data.InputParams += ",common.HexToAddress(ServConf.Common.TestAddr)"
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
	outfile, err := os.OpenFile(ServConf.Common.BuildPath+"/"+ServConf.Common.MainCodeName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}
	defer outfile.Close()
	// 读取abi文件信息
	abiInfos, err := readAbi(ServConf.Common.GoPath + "/" + ServConf.Common.ContractName + ".abi")
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

func Impl_config_code() error {
	//1. 写到哪
	outfile, err := os.OpenFile(ServConf.Common.BuildPath+"/"+ServConf.Common.ConfigCodeName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}
	defer outfile.Close()
	_, err = outfile.WriteString(config_build_temp)
	if err != nil {
		fmt.Println("failed to WriteString config", err)
		return err
	}
	//创建一个文件-config.toml
	outfile2, err := os.OpenFile(ServConf.Common.BuildPath+"/"+ServConf.Common.ConfigTomlName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("failed to open file", err)
		return err
	}
	defer outfile2.Close()
	//构建模版
	tmpl, err := template.New("config_toml").Parse(config_toml_temp)
	if err != nil {
		fmt.Println("failed to parse config_toml_temp", err)
		return err
	}
	data := ServConf.Version
	data.BuildDay = time.Now().Format("2006-01-02")
	err = tmpl.Execute(outfile2, data)
	if err != nil {
		fmt.Println("failed to Execute config_toml_temp", err)
		return err
	}
	return nil
}

func Run() {
	Impl_run_code()
	Impl_main_code()
	Impl_config_code()
}
