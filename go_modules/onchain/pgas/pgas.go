// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package pgas

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// Message is an auto generated low-level Go binding around an user-defined struct.
type Message struct {
	From  common.Address
	Nonce *big.Int
	Order *big.Int
	Start *big.Int
	To    common.Address
	Gas   *big.Int
	Data  []byte
}

// PGasMetaData contains all meta data concerning the PGas contract.
var PGasMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"order\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"messageValidate\",\"outputs\":[{\"internalType\":\"enumValidation\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PGasABI is the input ABI used to generate the binding from.
// Deprecated: Use PGasMetaData.ABI instead.
var PGasABI = PGasMetaData.ABI

// PGas is an auto generated Go binding around an Ethereum contract.
type PGas struct {
	PGasCaller     // Read-only binding to the contract
	PGasTransactor // Write-only binding to the contract
	PGasFilterer   // Log filterer for contract events
}

// PGasCaller is an auto generated read-only Go binding around an Ethereum contract.
type PGasCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PGasTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PGasTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PGasFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PGasFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PGasSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PGasSession struct {
	Contract     *PGas             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PGasCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PGasCallerSession struct {
	Contract *PGasCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// PGasTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PGasTransactorSession struct {
	Contract     *PGasTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PGasRaw is an auto generated low-level Go binding around an Ethereum contract.
type PGasRaw struct {
	Contract *PGas // Generic contract binding to access the raw methods on
}

// PGasCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PGasCallerRaw struct {
	Contract *PGasCaller // Generic read-only contract binding to access the raw methods on
}

// PGasTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PGasTransactorRaw struct {
	Contract *PGasTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPGas creates a new instance of PGas, bound to a specific deployed contract.
func NewPGas(address common.Address, backend bind.ContractBackend) (*PGas, error) {
	contract, err := bindPGas(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PGas{PGasCaller: PGasCaller{contract: contract}, PGasTransactor: PGasTransactor{contract: contract}, PGasFilterer: PGasFilterer{contract: contract}}, nil
}

// NewPGasCaller creates a new read-only instance of PGas, bound to a specific deployed contract.
func NewPGasCaller(address common.Address, caller bind.ContractCaller) (*PGasCaller, error) {
	contract, err := bindPGas(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PGasCaller{contract: contract}, nil
}

// NewPGasTransactor creates a new write-only instance of PGas, bound to a specific deployed contract.
func NewPGasTransactor(address common.Address, transactor bind.ContractTransactor) (*PGasTransactor, error) {
	contract, err := bindPGas(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PGasTransactor{contract: contract}, nil
}

// NewPGasFilterer creates a new log filterer instance of PGas, bound to a specific deployed contract.
func NewPGasFilterer(address common.Address, filterer bind.ContractFilterer) (*PGasFilterer, error) {
	contract, err := bindPGas(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PGasFilterer{contract: contract}, nil
}

// bindPGas binds a generic wrapper to an already deployed contract.
func bindPGas(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PGasMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PGas *PGasRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PGas.Contract.PGasCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PGas *PGasRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PGas.Contract.PGasTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PGas *PGasRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PGas.Contract.PGasTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PGas *PGasCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PGas.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PGas *PGasTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PGas.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PGas *PGasTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PGas.Contract.contract.Transact(opts, method, params...)
}

// MessageValidate is a free data retrieval call binding the contract method 0x98548a38.
//
// Solidity: function messageValidate((address,uint256,uint256,uint256,address,uint256,bytes) message) view returns(uint8)
func (_PGas *PGasCaller) MessageValidate(opts *bind.CallOpts, message Message) (uint8, error) {
	var out []interface{}
	err := _PGas.contract.Call(opts, &out, "messageValidate", message)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// MessageValidate is a free data retrieval call binding the contract method 0x98548a38.
//
// Solidity: function messageValidate((address,uint256,uint256,uint256,address,uint256,bytes) message) view returns(uint8)
func (_PGas *PGasSession) MessageValidate(message Message) (uint8, error) {
	return _PGas.Contract.MessageValidate(&_PGas.CallOpts, message)
}

// MessageValidate is a free data retrieval call binding the contract method 0x98548a38.
//
// Solidity: function messageValidate((address,uint256,uint256,uint256,address,uint256,bytes) message) view returns(uint8)
func (_PGas *PGasCallerSession) MessageValidate(message Message) (uint8, error) {
	return _PGas.Contract.MessageValidate(&_PGas.CallOpts, message)
}
