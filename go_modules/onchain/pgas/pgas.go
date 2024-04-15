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

// FilteredOrder is an auto generated low-level Go binding around an user-defined struct.
type FilteredOrder struct {
	Id       *big.Int
	Order    Order
	Status   uint8
	GasLeft  *big.Int
	Executor common.Address
}

// GasPayment is an auto generated low-level Go binding around an user-defined struct.
type GasPayment struct {
	Token   common.Address
	PerUnit *big.Int
}

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

// Order is an auto generated low-level Go binding around an user-defined struct.
type Order struct {
	Manager      common.Address
	Gas          *big.Int
	Expire       *big.Int
	Start        *big.Int
	End          *big.Int
	TxWindow     *big.Int
	RedeemWindow *big.Int
	GasPrice     GasPayment
	GasGuarantee GasPayment
}

// PGasMetaData contains all meta data concerning the PGas contract.
var PGasMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"domainSeparator\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"order\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structMessage\",\"name\":\"message\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"execute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"promisor\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"onlyLive\",\"type\":\"bool\"},{\"internalType\":\"uint256\",\"name\":\"limit\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"offset\",\"type\":\"uint256\"}],\"name\":\"getExecutorOrders\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"manager\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gas\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"expire\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"end\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"txWindow\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"redeemWindow\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"perUnit\",\"type\":\"uint256\"}],\"internalType\":\"structGasPayment\",\"name\":\"gasPrice\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"perUnit\",\"type\":\"uint256\"}],\"internalType\":\"structGasPayment\",\"name\":\"gasGuarantee\",\"type\":\"tuple\"}],\"internalType\":\"structOrder\",\"name\":\"order\",\"type\":\"tuple\"},{\"internalType\":\"enumOrderStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"gasLeft\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"executor\",\"type\":\"address\"}],\"internalType\":\"structFilteredOrder[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"order\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"start\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"gas\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"internalType\":\"structMessage\",\"name\":\"message\",\"type\":\"tuple\"}],\"name\":\"messageValidate\",\"outputs\":[{\"internalType\":\"enumValidation\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"nonce\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_PGas *PGasCaller) DomainSeparator(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _PGas.contract.Call(opts, &out, "domainSeparator")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_PGas *PGasSession) DomainSeparator() ([32]byte, error) {
	return _PGas.Contract.DomainSeparator(&_PGas.CallOpts)
}

// DomainSeparator is a free data retrieval call binding the contract method 0xf698da25.
//
// Solidity: function domainSeparator() view returns(bytes32)
func (_PGas *PGasCallerSession) DomainSeparator() ([32]byte, error) {
	return _PGas.Contract.DomainSeparator(&_PGas.CallOpts)
}

// GetExecutorOrders is a free data retrieval call binding the contract method 0x8d6aef3e.
//
// Solidity: function getExecutorOrders(address promisor, bool onlyLive, uint256 limit, uint256 offset) view returns((uint256,(address,uint256,uint256,uint256,uint256,uint256,uint256,(address,uint256),(address,uint256)),uint8,uint256,address)[])
func (_PGas *PGasCaller) GetExecutorOrders(opts *bind.CallOpts, promisor common.Address, onlyLive bool, limit *big.Int, offset *big.Int) ([]FilteredOrder, error) {
	var out []interface{}
	err := _PGas.contract.Call(opts, &out, "getExecutorOrders", promisor, onlyLive, limit, offset)

	if err != nil {
		return *new([]FilteredOrder), err
	}

	out0 := *abi.ConvertType(out[0], new([]FilteredOrder)).(*[]FilteredOrder)

	return out0, err

}

// GetExecutorOrders is a free data retrieval call binding the contract method 0x8d6aef3e.
//
// Solidity: function getExecutorOrders(address promisor, bool onlyLive, uint256 limit, uint256 offset) view returns((uint256,(address,uint256,uint256,uint256,uint256,uint256,uint256,(address,uint256),(address,uint256)),uint8,uint256,address)[])
func (_PGas *PGasSession) GetExecutorOrders(promisor common.Address, onlyLive bool, limit *big.Int, offset *big.Int) ([]FilteredOrder, error) {
	return _PGas.Contract.GetExecutorOrders(&_PGas.CallOpts, promisor, onlyLive, limit, offset)
}

// GetExecutorOrders is a free data retrieval call binding the contract method 0x8d6aef3e.
//
// Solidity: function getExecutorOrders(address promisor, bool onlyLive, uint256 limit, uint256 offset) view returns((uint256,(address,uint256,uint256,uint256,uint256,uint256,uint256,(address,uint256),(address,uint256)),uint8,uint256,address)[])
func (_PGas *PGasCallerSession) GetExecutorOrders(promisor common.Address, onlyLive bool, limit *big.Int, offset *big.Int) ([]FilteredOrder, error) {
	return _PGas.Contract.GetExecutorOrders(&_PGas.CallOpts, promisor, onlyLive, limit, offset)
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

// Nonce is a free data retrieval call binding the contract method 0xe8b1a9b2.
//
// Solidity: function nonce(address , uint256 ) view returns(bool)
func (_PGas *PGasCaller) Nonce(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _PGas.contract.Call(opts, &out, "nonce", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Nonce is a free data retrieval call binding the contract method 0xe8b1a9b2.
//
// Solidity: function nonce(address , uint256 ) view returns(bool)
func (_PGas *PGasSession) Nonce(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _PGas.Contract.Nonce(&_PGas.CallOpts, arg0, arg1)
}

// Nonce is a free data retrieval call binding the contract method 0xe8b1a9b2.
//
// Solidity: function nonce(address , uint256 ) view returns(bool)
func (_PGas *PGasCallerSession) Nonce(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _PGas.Contract.Nonce(&_PGas.CallOpts, arg0, arg1)
}

// Execute is a paid mutator transaction binding the contract method 0x03f6a219.
//
// Solidity: function execute((address,uint256,uint256,uint256,address,uint256,bytes) message, bytes signature) returns()
func (_PGas *PGasTransactor) Execute(opts *bind.TransactOpts, message Message, signature []byte) (*types.Transaction, error) {
	return _PGas.contract.Transact(opts, "execute", message, signature)
}

// Execute is a paid mutator transaction binding the contract method 0x03f6a219.
//
// Solidity: function execute((address,uint256,uint256,uint256,address,uint256,bytes) message, bytes signature) returns()
func (_PGas *PGasSession) Execute(message Message, signature []byte) (*types.Transaction, error) {
	return _PGas.Contract.Execute(&_PGas.TransactOpts, message, signature)
}

// Execute is a paid mutator transaction binding the contract method 0x03f6a219.
//
// Solidity: function execute((address,uint256,uint256,uint256,address,uint256,bytes) message, bytes signature) returns()
func (_PGas *PGasTransactorSession) Execute(message Message, signature []byte) (*types.Transaction, error) {
	return _PGas.Contract.Execute(&_PGas.TransactOpts, message, signature)
}
