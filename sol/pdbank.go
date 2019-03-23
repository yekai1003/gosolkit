// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PdbankABI is the input ABI used to generate the binding from.
const PdbankABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"totalAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bankName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_bankName\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// PdbankBin is the compiled bytecode used for deploying new contracts.
const PdbankBin = `0x608060405234801561001057600080fd5b506040516104533803806104538339810180604052602081101561003357600080fd5b81019080805164010000000081111561004b57600080fd5b8201602081018481111561005e57600080fd5b815164010000000081118282018710171561007857600080fd5b5050600080546001600160a01b0319163317905580519093506100a492506003915060208401906100ab565b5050610146565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100ec57805160ff1916838001178555610119565b82800160010185558215610119579182015b828111156101195782518255916020019190600101906100fe565b50610125929150610129565b5090565b61014391905b80821115610125576000815560010161012f565b90565b6102fe806101556000396000f3fe6080604052600436106100555760003560e01c80631a39d8ef1461005a57806327d358011461008157806327e235e31461010b5780632e1a7d4d1461013e5780638da5cb5b1461015d578063d0e30db01461018e575b600080fd5b34801561006657600080fd5b5061006f610196565b60408051918252519081900360200190f35b34801561008d57600080fd5b5061009661019c565b6040805160208082528351818301528351919283929083019185019080838360005b838110156100d05781810151838201526020016100b8565b50505050905090810190601f1680156100fd5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561011757600080fd5b5061006f6004803603602081101561012e57600080fd5b50356001600160a01b031661022a565b61015b6004803603602081101561015457600080fd5b503561023c565b005b34801561016957600080fd5b506101726102a0565b604080516001600160a01b039092168252519081900360200190f35b61015b6102af565b60025481565b6003805460408051602060026001851615610100026000190190941693909304601f810184900484028201840190925281815292918301828280156102225780601f106101f757610100808354040283529160200191610222565b820191906000526020600020905b81548152906001019060200180831161020557829003601f168201915b505050505081565b60016020526000908152604090205481565b3360009081526001602052604090205481101561029d5733600081815260016020526040808220805485900390555183156108fc0291849190818181858888f19350505050158015610292573d6000803e3d6000fd5b506002805482900390555b50565b6000546001600160a01b031681565b60028054349081019091553360009081526001602052604090208054909101905556fea165627a7a72305820629748dc1f1653dff1bd207f2f4189417300c71c2614eed1aa6ee70a567a92710029`

// DeployPdbank deploys a new Ethereum contract, binding an instance of Pdbank to it.
func DeployPdbank(auth *bind.TransactOpts, backend bind.ContractBackend, _bankName string) (common.Address, *types.Transaction, *Pdbank, error) {
	parsed, err := abi.JSON(strings.NewReader(PdbankABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PdbankBin), backend, _bankName)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Pdbank{PdbankCaller: PdbankCaller{contract: contract}, PdbankTransactor: PdbankTransactor{contract: contract}, PdbankFilterer: PdbankFilterer{contract: contract}}, nil
}

// Pdbank is an auto generated Go binding around an Ethereum contract.
type Pdbank struct {
	PdbankCaller     // Read-only binding to the contract
	PdbankTransactor // Write-only binding to the contract
	PdbankFilterer   // Log filterer for contract events
}

// PdbankCaller is an auto generated read-only Go binding around an Ethereum contract.
type PdbankCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PdbankTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PdbankTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PdbankFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PdbankFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PdbankSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PdbankSession struct {
	Contract     *Pdbank           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PdbankCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PdbankCallerSession struct {
	Contract *PdbankCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PdbankTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PdbankTransactorSession struct {
	Contract     *PdbankTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PdbankRaw is an auto generated low-level Go binding around an Ethereum contract.
type PdbankRaw struct {
	Contract *Pdbank // Generic contract binding to access the raw methods on
}

// PdbankCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PdbankCallerRaw struct {
	Contract *PdbankCaller // Generic read-only contract binding to access the raw methods on
}

// PdbankTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PdbankTransactorRaw struct {
	Contract *PdbankTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPdbank creates a new instance of Pdbank, bound to a specific deployed contract.
func NewPdbank(address common.Address, backend bind.ContractBackend) (*Pdbank, error) {
	contract, err := bindPdbank(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Pdbank{PdbankCaller: PdbankCaller{contract: contract}, PdbankTransactor: PdbankTransactor{contract: contract}, PdbankFilterer: PdbankFilterer{contract: contract}}, nil
}

// NewPdbankCaller creates a new read-only instance of Pdbank, bound to a specific deployed contract.
func NewPdbankCaller(address common.Address, caller bind.ContractCaller) (*PdbankCaller, error) {
	contract, err := bindPdbank(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PdbankCaller{contract: contract}, nil
}

// NewPdbankTransactor creates a new write-only instance of Pdbank, bound to a specific deployed contract.
func NewPdbankTransactor(address common.Address, transactor bind.ContractTransactor) (*PdbankTransactor, error) {
	contract, err := bindPdbank(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PdbankTransactor{contract: contract}, nil
}

// NewPdbankFilterer creates a new log filterer instance of Pdbank, bound to a specific deployed contract.
func NewPdbankFilterer(address common.Address, filterer bind.ContractFilterer) (*PdbankFilterer, error) {
	contract, err := bindPdbank(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PdbankFilterer{contract: contract}, nil
}

// bindPdbank binds a generic wrapper to an already deployed contract.
func bindPdbank(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PdbankABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pdbank *PdbankRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Pdbank.Contract.PdbankCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pdbank *PdbankRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pdbank.Contract.PdbankTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pdbank *PdbankRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pdbank.Contract.PdbankTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Pdbank *PdbankCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Pdbank.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Pdbank *PdbankTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pdbank.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Pdbank *PdbankTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Pdbank.Contract.contract.Transact(opts, method, params...)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_Pdbank *PdbankCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Pdbank.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_Pdbank *PdbankSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _Pdbank.Contract.Balances(&_Pdbank.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_Pdbank *PdbankCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _Pdbank.Contract.Balances(&_Pdbank.CallOpts, arg0)
}

// BankName is a free data retrieval call binding the contract method 0x27d35801.
//
// Solidity: function bankName() constant returns(string)
func (_Pdbank *PdbankCaller) BankName(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Pdbank.contract.Call(opts, out, "bankName")
	return *ret0, err
}

// BankName is a free data retrieval call binding the contract method 0x27d35801.
//
// Solidity: function bankName() constant returns(string)
func (_Pdbank *PdbankSession) BankName() (string, error) {
	return _Pdbank.Contract.BankName(&_Pdbank.CallOpts)
}

// BankName is a free data retrieval call binding the contract method 0x27d35801.
//
// Solidity: function bankName() constant returns(string)
func (_Pdbank *PdbankCallerSession) BankName() (string, error) {
	return _Pdbank.Contract.BankName(&_Pdbank.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Pdbank *PdbankCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Pdbank.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Pdbank *PdbankSession) Owner() (common.Address, error) {
	return _Pdbank.Contract.Owner(&_Pdbank.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Pdbank *PdbankCallerSession) Owner() (common.Address, error) {
	return _Pdbank.Contract.Owner(&_Pdbank.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() constant returns(uint256)
func (_Pdbank *PdbankCaller) TotalAmount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Pdbank.contract.Call(opts, out, "totalAmount")
	return *ret0, err
}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() constant returns(uint256)
func (_Pdbank *PdbankSession) TotalAmount() (*big.Int, error) {
	return _Pdbank.Contract.TotalAmount(&_Pdbank.CallOpts)
}

// TotalAmount is a free data retrieval call binding the contract method 0x1a39d8ef.
//
// Solidity: function totalAmount() constant returns(uint256)
func (_Pdbank *PdbankCallerSession) TotalAmount() (*big.Int, error) {
	return _Pdbank.Contract.TotalAmount(&_Pdbank.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_Pdbank *PdbankTransactor) Deposit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Pdbank.contract.Transact(opts, "deposit")
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_Pdbank *PdbankSession) Deposit() (*types.Transaction, error) {
	return _Pdbank.Contract.Deposit(&_Pdbank.TransactOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xd0e30db0.
//
// Solidity: function deposit() returns()
func (_Pdbank *PdbankTransactorSession) Deposit() (*types.Transaction, error) {
	return _Pdbank.Contract.Deposit(&_Pdbank.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_Pdbank *PdbankTransactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _Pdbank.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_Pdbank *PdbankSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _Pdbank.Contract.Withdraw(&_Pdbank.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_Pdbank *PdbankTransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _Pdbank.Contract.Withdraw(&_Pdbank.TransactOpts, _amount)
}
