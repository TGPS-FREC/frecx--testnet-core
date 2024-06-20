// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"strings"

	"github.com/FRECNET/accounts/abi"
	"github.com/FRECNET/accounts/abi/bind"
	"github.com/FRECNET/common"
	"github.com/FRECNET/core/types"
)

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
const SafeMathBin = `0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146060604052600080fd00a165627a7a72305820b9407d48ebc7efee5c9f08b3b3a957df2939281f5913225e8c1291f069b900490029`

// DeploySafeMath deploys a new Ethereum contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Ethereum contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Ethereum contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Ethereum contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}

// FRERandomizeABI is the input ABI used to generate the binding from.
const FRERandomizeABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"getSecret\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_secret\",\"type\":\"bytes32[]\"}],\"name\":\"setSecret\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"getOpening\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_opening\",\"type\":\"bytes32\"}],\"name\":\"setOpening\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// FRERandomizeBin is the compiled bytecode used for deploying new contracts.
const FRERandomizeBin = `0x6060604052341561000f57600080fd5b6103368061001e6000396000f3006060604052600436106100615763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663284180fc811461006657806334d38600146100d8578063d442d6cc14610129578063e11f5ba21461015a575b600080fd5b341561007157600080fd5b610085600160a060020a0360043516610170565b60405160208082528190810183818151815260200191508051906020019060200280838360005b838110156100c45780820151838201526020016100ac565b505050509050019250505060405180910390f35b34156100e357600080fd5b61012760046024813581810190830135806020818102016040519081016040528093929190818152602001838360200280828437509496506101f395505050505050565b005b341561013457600080fd5b610148600160a060020a0360043516610243565b60405190815260200160405180910390f35b341561016557600080fd5b61012760043561025e565b61017861028e565b60008083600160a060020a0316600160a060020a031681526020019081526020016000208054806020026020016040519081016040528092919081815260200182805480156101e757602002820191906000526020600020905b815481526001909101906020018083116101d2575b50505050509050919050565b610384430661032081101561020757600080fd5b610352811061021557600080fd5b600160a060020a033316600090815260208190526040902082805161023e9291602001906102a0565b505050565b600160a060020a031660009081526001602052604090205490565b610384430661035281101561027257600080fd5b50600160a060020a033316600090815260016020526040902055565b60206040519081016040526000815290565b8280548282559060005260206000209081019282156102dd579160200282015b828111156102dd57825182556020909201916001909101906102c0565b506102e99291506102ed565b5090565b61030791905b808211156102e957600081556001016102f3565b905600a165627a7a7230582034991c8dc4001fc254f3ba2811c05d2e7d29bee3908946ca56d1545b2c852de20029`

// DeployFRERandomize deploys a new Ethereum contract, binding an instance of FRERandomize to it.
func DeployFRERandomize(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *FRERandomize, error) {
	parsed, err := abi.JSON(strings.NewReader(FRERandomizeABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FRERandomizeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FRERandomize{FRERandomizeCaller: FRERandomizeCaller{contract: contract}, FRERandomizeTransactor: FRERandomizeTransactor{contract: contract}, FRERandomizeFilterer: FRERandomizeFilterer{contract: contract}}, nil
}

// FRERandomize is an auto generated Go binding around an Ethereum contract.
type FRERandomize struct {
	FRERandomizeCaller     // Read-only binding to the contract
	FRERandomizeTransactor // Write-only binding to the contract
	FRERandomizeFilterer   // Log filterer for contract events
}

// FRERandomizeCaller is an auto generated read-only Go binding around an Ethereum contract.
type FRERandomizeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FRERandomizeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FRERandomizeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FRERandomizeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FRERandomizeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FRERandomizeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FRERandomizeSession struct {
	Contract     *FRERandomize     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FRERandomizeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FRERandomizeCallerSession struct {
	Contract *FRERandomizeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// FRERandomizeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FRERandomizeTransactorSession struct {
	Contract     *FRERandomizeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// FRERandomizeRaw is an auto generated low-level Go binding around an Ethereum contract.
type FRERandomizeRaw struct {
	Contract *FRERandomize // Generic contract binding to access the raw methods on
}

// FRERandomizeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FRERandomizeCallerRaw struct {
	Contract *FRERandomizeCaller // Generic read-only contract binding to access the raw methods on
}

// FRERandomizeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FRERandomizeTransactorRaw struct {
	Contract *FRERandomizeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFRERandomize creates a new instance of FRERandomize, bound to a specific deployed contract.
func NewFRERandomize(address common.Address, backend bind.ContractBackend) (*FRERandomize, error) {
	contract, err := bindFRERandomize(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FRERandomize{FRERandomizeCaller: FRERandomizeCaller{contract: contract}, FRERandomizeTransactor: FRERandomizeTransactor{contract: contract}, FRERandomizeFilterer: FRERandomizeFilterer{contract: contract}}, nil
}

// NewFRERandomizeCaller creates a new read-only instance of FRERandomize, bound to a specific deployed contract.
func NewFRERandomizeCaller(address common.Address, caller bind.ContractCaller) (*FRERandomizeCaller, error) {
	contract, err := bindFRERandomize(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FRERandomizeCaller{contract: contract}, nil
}

// NewFRERandomizeTransactor creates a new write-only instance of FRERandomize, bound to a specific deployed contract.
func NewFRERandomizeTransactor(address common.Address, transactor bind.ContractTransactor) (*FRERandomizeTransactor, error) {
	contract, err := bindFRERandomize(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FRERandomizeTransactor{contract: contract}, nil
}

// NewFRERandomizeFilterer creates a new log filterer instance of FRERandomize, bound to a specific deployed contract.
func NewFRERandomizeFilterer(address common.Address, filterer bind.ContractFilterer) (*FRERandomizeFilterer, error) {
	contract, err := bindFRERandomize(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FRERandomizeFilterer{contract: contract}, nil
}

// bindFRERandomize binds a generic wrapper to an already deployed contract.
func bindFRERandomize(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FRERandomizeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FRERandomize *FRERandomizeRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FRERandomize.Contract.FRERandomizeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FRERandomize *FRERandomizeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FRERandomize.Contract.FRERandomizeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FRERandomize *FRERandomizeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FRERandomize.Contract.FRERandomizeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FRERandomize *FRERandomizeCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FRERandomize.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FRERandomize *FRERandomizeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FRERandomize.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FRERandomize *FRERandomizeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FRERandomize.Contract.contract.Transact(opts, method, params...)
}

// GetOpening is a free data retrieval call binding the contract method 0xd442d6cc.
//
// Solidity: function getOpening(_validator address) constant returns(bytes32)
func (_FRERandomize *FRERandomizeCaller) GetOpening(opts *bind.CallOpts, _validator common.Address) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _FRERandomize.contract.Call(opts, out, "getOpening", _validator)
	return *ret0, err
}

// GetOpening is a free data retrieval call binding the contract method 0xd442d6cc.
//
// Solidity: function getOpening(_validator address) constant returns(bytes32)
func (_FRERandomize *FRERandomizeSession) GetOpening(_validator common.Address) ([32]byte, error) {
	return _FRERandomize.Contract.GetOpening(&_FRERandomize.CallOpts, _validator)
}

// GetOpening is a free data retrieval call binding the contract method 0xd442d6cc.
//
// Solidity: function getOpening(_validator address) constant returns(bytes32)
func (_FRERandomize *FRERandomizeCallerSession) GetOpening(_validator common.Address) ([32]byte, error) {
	return _FRERandomize.Contract.GetOpening(&_FRERandomize.CallOpts, _validator)
}

// GetSecret is a free data retrieval call binding the contract method 0x284180fc.
//
// Solidity: function getSecret(_validator address) constant returns(bytes32[])
func (_FRERandomize *FRERandomizeCaller) GetSecret(opts *bind.CallOpts, _validator common.Address) ([][32]byte, error) {
	var (
		ret0 = new([][32]byte)
	)
	out := ret0
	err := _FRERandomize.contract.Call(opts, out, "getSecret", _validator)
	return *ret0, err
}

// GetSecret is a free data retrieval call binding the contract method 0x284180fc.
//
// Solidity: function getSecret(_validator address) constant returns(bytes32[])
func (_FRERandomize *FRERandomizeSession) GetSecret(_validator common.Address) ([][32]byte, error) {
	return _FRERandomize.Contract.GetSecret(&_FRERandomize.CallOpts, _validator)
}

// GetSecret is a free data retrieval call binding the contract method 0x284180fc.
//
// Solidity: function getSecret(_validator address) constant returns(bytes32[])
func (_FRERandomize *FRERandomizeCallerSession) GetSecret(_validator common.Address) ([][32]byte, error) {
	return _FRERandomize.Contract.GetSecret(&_FRERandomize.CallOpts, _validator)
}

// SetOpening is a paid mutator transaction binding the contract method 0xe11f5ba2.
//
// Solidity: function setOpening(_opening bytes32) returns()
func (_FRERandomize *FRERandomizeTransactor) SetOpening(opts *bind.TransactOpts, _opening [32]byte) (*types.Transaction, error) {
	return _FRERandomize.contract.Transact(opts, "setOpening", _opening)
}

// SetOpening is a paid mutator transaction binding the contract method 0xe11f5ba2.
//
// Solidity: function setOpening(_opening bytes32) returns()
func (_FRERandomize *FRERandomizeSession) SetOpening(_opening [32]byte) (*types.Transaction, error) {
	return _FRERandomize.Contract.SetOpening(&_FRERandomize.TransactOpts, _opening)
}

// SetOpening is a paid mutator transaction binding the contract method 0xe11f5ba2.
//
// Solidity: function setOpening(_opening bytes32) returns()
func (_FRERandomize *FRERandomizeTransactorSession) SetOpening(_opening [32]byte) (*types.Transaction, error) {
	return _FRERandomize.Contract.SetOpening(&_FRERandomize.TransactOpts, _opening)
}

// SetSecret is a paid mutator transaction binding the contract method 0x34d38600.
//
// Solidity: function setSecret(_secret bytes32[]) returns()
func (_FRERandomize *FRERandomizeTransactor) SetSecret(opts *bind.TransactOpts, _secret [][32]byte) (*types.Transaction, error) {
	return _FRERandomize.contract.Transact(opts, "setSecret", _secret)
}

// SetSecret is a paid mutator transaction binding the contract method 0x34d38600.
//
// Solidity: function setSecret(_secret bytes32[]) returns()
func (_FRERandomize *FRERandomizeSession) SetSecret(_secret [][32]byte) (*types.Transaction, error) {
	return _FRERandomize.Contract.SetSecret(&_FRERandomize.TransactOpts, _secret)
}

// SetSecret is a paid mutator transaction binding the contract method 0x34d38600.
//
// Solidity: function setSecret(_secret bytes32[]) returns()
func (_FRERandomize *FRERandomizeTransactorSession) SetSecret(_secret [][32]byte) (*types.Transaction, error) {
	return _FRERandomize.Contract.SetSecret(&_FRERandomize.TransactOpts, _secret)
}
