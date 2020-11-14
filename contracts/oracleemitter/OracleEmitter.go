// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package oracleEmitter

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
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OracleEmitterABI is the input ABI used to generate the binding from.
const OracleEmitterABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"originChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestParams\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"callbackMethodID\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"NewOracleRequest\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"originChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"requestType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"requestParams\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackMethodID\",\"type\":\"bytes4\"}],\"name\":\"requestOracles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// OracleEmitter is an auto generated Go binding around an Ethereum contract.
type OracleEmitter struct {
	OracleEmitterCaller     // Read-only binding to the contract
	OracleEmitterTransactor // Write-only binding to the contract
	OracleEmitterFilterer   // Log filterer for contract events
}

// OracleEmitterCaller is an auto generated read-only Go binding around an Ethereum contract.
type OracleEmitterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleEmitterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OracleEmitterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleEmitterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OracleEmitterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OracleEmitterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OracleEmitterSession struct {
	Contract     *OracleEmitter    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OracleEmitterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OracleEmitterCallerSession struct {
	Contract *OracleEmitterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// OracleEmitterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OracleEmitterTransactorSession struct {
	Contract     *OracleEmitterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// OracleEmitterRaw is an auto generated low-level Go binding around an Ethereum contract.
type OracleEmitterRaw struct {
	Contract *OracleEmitter // Generic contract binding to access the raw methods on
}

// OracleEmitterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OracleEmitterCallerRaw struct {
	Contract *OracleEmitterCaller // Generic read-only contract binding to access the raw methods on
}

// OracleEmitterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OracleEmitterTransactorRaw struct {
	Contract *OracleEmitterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOracleEmitter creates a new instance of OracleEmitter, bound to a specific deployed contract.
func NewOracleEmitter(address common.Address, backend bind.ContractBackend) (*OracleEmitter, error) {
	contract, err := bindOracleEmitter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OracleEmitter{OracleEmitterCaller: OracleEmitterCaller{contract: contract}, OracleEmitterTransactor: OracleEmitterTransactor{contract: contract}, OracleEmitterFilterer: OracleEmitterFilterer{contract: contract}}, nil
}

// NewOracleEmitterCaller creates a new read-only instance of OracleEmitter, bound to a specific deployed contract.
func NewOracleEmitterCaller(address common.Address, caller bind.ContractCaller) (*OracleEmitterCaller, error) {
	contract, err := bindOracleEmitter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OracleEmitterCaller{contract: contract}, nil
}

// NewOracleEmitterTransactor creates a new write-only instance of OracleEmitter, bound to a specific deployed contract.
func NewOracleEmitterTransactor(address common.Address, transactor bind.ContractTransactor) (*OracleEmitterTransactor, error) {
	contract, err := bindOracleEmitter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OracleEmitterTransactor{contract: contract}, nil
}

// NewOracleEmitterFilterer creates a new log filterer instance of OracleEmitter, bound to a specific deployed contract.
func NewOracleEmitterFilterer(address common.Address, filterer bind.ContractFilterer) (*OracleEmitterFilterer, error) {
	contract, err := bindOracleEmitter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OracleEmitterFilterer{contract: contract}, nil
}

// bindOracleEmitter binds a generic wrapper to an already deployed contract.
func bindOracleEmitter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OracleEmitterABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleEmitter *OracleEmitterRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OracleEmitter.Contract.OracleEmitterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleEmitter *OracleEmitterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleEmitter.Contract.OracleEmitterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleEmitter *OracleEmitterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleEmitter.Contract.OracleEmitterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OracleEmitter *OracleEmitterCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OracleEmitter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OracleEmitter *OracleEmitterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OracleEmitter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OracleEmitter *OracleEmitterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OracleEmitter.Contract.contract.Transact(opts, method, params...)
}

// RequestOracles is a paid mutator transaction binding the contract method 0x8fe3d79f.
//
// Solidity: function requestOracles(string originChain, string requestType, string requestParams, address callbackAddress, bytes4 callbackMethodID) returns(uint256)
func (_OracleEmitter *OracleEmitterTransactor) RequestOracles(opts *bind.TransactOpts, originChain string, requestType string, requestParams string, callbackAddress common.Address, callbackMethodID [4]byte) (*types.Transaction, error) {
	return _OracleEmitter.contract.Transact(opts, "requestOracles", originChain, requestType, requestParams, callbackAddress, callbackMethodID)
}

// RequestOracles is a paid mutator transaction binding the contract method 0x8fe3d79f.
//
// Solidity: function requestOracles(string originChain, string requestType, string requestParams, address callbackAddress, bytes4 callbackMethodID) returns(uint256)
func (_OracleEmitter *OracleEmitterSession) RequestOracles(originChain string, requestType string, requestParams string, callbackAddress common.Address, callbackMethodID [4]byte) (*types.Transaction, error) {
	return _OracleEmitter.Contract.RequestOracles(&_OracleEmitter.TransactOpts, originChain, requestType, requestParams, callbackAddress, callbackMethodID)
}

// RequestOracles is a paid mutator transaction binding the contract method 0x8fe3d79f.
//
// Solidity: function requestOracles(string originChain, string requestType, string requestParams, address callbackAddress, bytes4 callbackMethodID) returns(uint256)
func (_OracleEmitter *OracleEmitterTransactorSession) RequestOracles(originChain string, requestType string, requestParams string, callbackAddress common.Address, callbackMethodID [4]byte) (*types.Transaction, error) {
	return _OracleEmitter.Contract.RequestOracles(&_OracleEmitter.TransactOpts, originChain, requestType, requestParams, callbackAddress, callbackMethodID)
}

// OracleEmitterNewOracleRequestIterator is returned from FilterNewOracleRequest and is used to iterate over the raw logs and unpacked data for NewOracleRequest events raised by the OracleEmitter contract.
type OracleEmitterNewOracleRequestIterator struct {
	Event *OracleEmitterNewOracleRequest // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *OracleEmitterNewOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(OracleEmitterNewOracleRequest)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(OracleEmitterNewOracleRequest)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *OracleEmitterNewOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *OracleEmitterNewOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// OracleEmitterNewOracleRequest represents a NewOracleRequest event raised by the OracleEmitter contract.
type OracleEmitterNewOracleRequest struct {
	OriginChain      string
	RequestType      string
	RequestParams    string
	CallbackAddress  common.Address
	CallbackMethodID [4]byte
	RequestID        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNewOracleRequest is a free log retrieval operation binding the contract event 0x4840d7d041230e7d05eaf9e64f0924379c740cd38ebe8ced4ca8b22afb8d3c94.
//
// Solidity: event NewOracleRequest(string originChain, string requestType, string requestParams, address callbackAddress, bytes4 callbackMethodID, uint256 requestID)
func (_OracleEmitter *OracleEmitterFilterer) FilterNewOracleRequest(opts *bind.FilterOpts) (*OracleEmitterNewOracleRequestIterator, error) {

	logs, sub, err := _OracleEmitter.contract.FilterLogs(opts, "NewOracleRequest")
	if err != nil {
		return nil, err
	}
	return &OracleEmitterNewOracleRequestIterator{contract: _OracleEmitter.contract, event: "NewOracleRequest", logs: logs, sub: sub}, nil
}

// WatchNewOracleRequest is a free log subscription operation binding the contract event 0x4840d7d041230e7d05eaf9e64f0924379c740cd38ebe8ced4ca8b22afb8d3c94.
//
// Solidity: event NewOracleRequest(string originChain, string requestType, string requestParams, address callbackAddress, bytes4 callbackMethodID, uint256 requestID)
func (_OracleEmitter *OracleEmitterFilterer) WatchNewOracleRequest(opts *bind.WatchOpts, sink chan<- *OracleEmitterNewOracleRequest) (event.Subscription, error) {

	logs, sub, err := _OracleEmitter.contract.WatchLogs(opts, "NewOracleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(OracleEmitterNewOracleRequest)
				if err := _OracleEmitter.contract.UnpackLog(event, "NewOracleRequest", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseNewOracleRequest is a log parse operation binding the contract event 0x4840d7d041230e7d05eaf9e64f0924379c740cd38ebe8ced4ca8b22afb8d3c94.
//
// Solidity: event NewOracleRequest(string originChain, string requestType, string requestParams, address callbackAddress, bytes4 callbackMethodID, uint256 requestID)
func (_OracleEmitter *OracleEmitterFilterer) ParseNewOracleRequest(log types.Log) (*OracleEmitterNewOracleRequest, error) {
	event := new(OracleEmitterNewOracleRequest)
	if err := _OracleEmitter.contract.UnpackLog(event, "NewOracleRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}
