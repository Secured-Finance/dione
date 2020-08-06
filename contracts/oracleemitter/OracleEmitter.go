// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package oracleemitter

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

// SmartcontractsABI is the input ABI used to generate the binding from.
const SmartcontractsABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes4\",\"name\":\"callbackMethodID\",\"type\":\"bytes4\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"NewOracleRequest\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"requestType\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"callbackMethodID\",\"type\":\"bytes4\"}],\"name\":\"requestOracles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Smartcontracts is an auto generated Go binding around an Ethereum contract.
type Smartcontracts struct {
	SmartcontractsCaller     // Read-only binding to the contract
	SmartcontractsTransactor // Write-only binding to the contract
	SmartcontractsFilterer   // Log filterer for contract events
}

// SmartcontractsCaller is an auto generated read-only Go binding around an Ethereum contract.
type SmartcontractsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmartcontractsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SmartcontractsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmartcontractsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SmartcontractsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SmartcontractsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SmartcontractsSession struct {
	Contract     *Smartcontracts   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SmartcontractsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SmartcontractsCallerSession struct {
	Contract *SmartcontractsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// SmartcontractsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SmartcontractsTransactorSession struct {
	Contract     *SmartcontractsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// SmartcontractsRaw is an auto generated low-level Go binding around an Ethereum contract.
type SmartcontractsRaw struct {
	Contract *Smartcontracts // Generic contract binding to access the raw methods on
}

// SmartcontractsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SmartcontractsCallerRaw struct {
	Contract *SmartcontractsCaller // Generic read-only contract binding to access the raw methods on
}

// SmartcontractsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SmartcontractsTransactorRaw struct {
	Contract *SmartcontractsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSmartcontracts creates a new instance of Smartcontracts, bound to a specific deployed contract.
func NewSmartcontracts(address common.Address, backend bind.ContractBackend) (*Smartcontracts, error) {
	contract, err := bindSmartcontracts(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Smartcontracts{SmartcontractsCaller: SmartcontractsCaller{contract: contract}, SmartcontractsTransactor: SmartcontractsTransactor{contract: contract}, SmartcontractsFilterer: SmartcontractsFilterer{contract: contract}}, nil
}

// NewSmartcontractsCaller creates a new read-only instance of Smartcontracts, bound to a specific deployed contract.
func NewSmartcontractsCaller(address common.Address, caller bind.ContractCaller) (*SmartcontractsCaller, error) {
	contract, err := bindSmartcontracts(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SmartcontractsCaller{contract: contract}, nil
}

// NewSmartcontractsTransactor creates a new write-only instance of Smartcontracts, bound to a specific deployed contract.
func NewSmartcontractsTransactor(address common.Address, transactor bind.ContractTransactor) (*SmartcontractsTransactor, error) {
	contract, err := bindSmartcontracts(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SmartcontractsTransactor{contract: contract}, nil
}

// NewSmartcontractsFilterer creates a new log filterer instance of Smartcontracts, bound to a specific deployed contract.
func NewSmartcontractsFilterer(address common.Address, filterer bind.ContractFilterer) (*SmartcontractsFilterer, error) {
	contract, err := bindSmartcontracts(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SmartcontractsFilterer{contract: contract}, nil
}

// bindSmartcontracts binds a generic wrapper to an already deployed contract.
func bindSmartcontracts(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SmartcontractsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Smartcontracts *SmartcontractsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Smartcontracts.Contract.SmartcontractsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Smartcontracts *SmartcontractsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Smartcontracts.Contract.SmartcontractsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Smartcontracts *SmartcontractsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Smartcontracts.Contract.SmartcontractsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Smartcontracts *SmartcontractsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Smartcontracts.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Smartcontracts *SmartcontractsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Smartcontracts.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Smartcontracts *SmartcontractsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Smartcontracts.Contract.contract.Transact(opts, method, params...)
}

// RequestOracles is a paid mutator transaction binding the contract method 0xcae2f480.
//
// Solidity: function requestOracles(string requestType, address callbackAddress, bytes4 callbackMethodID) returns(uint256)
func (_Smartcontracts *SmartcontractsTransactor) RequestOracles(opts *bind.TransactOpts, requestType string, callbackAddress common.Address, callbackMethodID [4]byte) (*types.Transaction, error) {
	return _Smartcontracts.contract.Transact(opts, "requestOracles", requestType, callbackAddress, callbackMethodID)
}

// RequestOracles is a paid mutator transaction binding the contract method 0xcae2f480.
//
// Solidity: function requestOracles(string requestType, address callbackAddress, bytes4 callbackMethodID) returns(uint256)
func (_Smartcontracts *SmartcontractsSession) RequestOracles(requestType string, callbackAddress common.Address, callbackMethodID [4]byte) (*types.Transaction, error) {
	return _Smartcontracts.Contract.RequestOracles(&_Smartcontracts.TransactOpts, requestType, callbackAddress, callbackMethodID)
}

// RequestOracles is a paid mutator transaction binding the contract method 0xcae2f480.
//
// Solidity: function requestOracles(string requestType, address callbackAddress, bytes4 callbackMethodID) returns(uint256)
func (_Smartcontracts *SmartcontractsTransactorSession) RequestOracles(requestType string, callbackAddress common.Address, callbackMethodID [4]byte) (*types.Transaction, error) {
	return _Smartcontracts.Contract.RequestOracles(&_Smartcontracts.TransactOpts, requestType, callbackAddress, callbackMethodID)
}

// SmartcontractsNewOracleRequestIterator is returned from FilterNewOracleRequest and is used to iterate over the raw logs and unpacked data for NewOracleRequest events raised by the Smartcontracts contract.
type SmartcontractsNewOracleRequestIterator struct {
	Event *SmartcontractsNewOracleRequest // Event containing the contract specifics and raw log

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
func (it *SmartcontractsNewOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmartcontractsNewOracleRequest)
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
		it.Event = new(SmartcontractsNewOracleRequest)
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
func (it *SmartcontractsNewOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmartcontractsNewOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmartcontractsNewOracleRequest represents a NewOracleRequest event raised by the Smartcontracts contract.
type SmartcontractsNewOracleRequest struct {
	RequestType      string
	CallbackAddress  common.Address
	CallbackMethodID [4]byte
	RequestID        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNewOracleRequest is a free log retrieval operation binding the contract event 0xc1a7df69ed0404441720a3eb710f322985ecbbe4f31934587bee703efe04427a.
//
// Solidity: event NewOracleRequest(string requestType, address callbackAddress, bytes4 callbackMethodID, uint256 requestID)
func (_Smartcontracts *SmartcontractsFilterer) FilterNewOracleRequest(opts *bind.FilterOpts) (*SmartcontractsNewOracleRequestIterator, error) {

	logs, sub, err := _Smartcontracts.contract.FilterLogs(opts, "NewOracleRequest")
	if err != nil {
		return nil, err
	}
	return &SmartcontractsNewOracleRequestIterator{contract: _Smartcontracts.contract, event: "NewOracleRequest", logs: logs, sub: sub}, nil
}

// WatchNewOracleRequest is a free log subscription operation binding the contract event 0xc1a7df69ed0404441720a3eb710f322985ecbbe4f31934587bee703efe04427a.
//
// Solidity: event NewOracleRequest(string requestType, address callbackAddress, bytes4 callbackMethodID, uint256 requestID)
func (_Smartcontracts *SmartcontractsFilterer) WatchNewOracleRequest(opts *bind.WatchOpts, sink chan<- *SmartcontractsNewOracleRequest) (event.Subscription, error) {

	logs, sub, err := _Smartcontracts.contract.WatchLogs(opts, "NewOracleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmartcontractsNewOracleRequest)
				if err := _Smartcontracts.contract.UnpackLog(event, "NewOracleRequest", log); err != nil {
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

// ParseNewOracleRequest is a log parse operation binding the contract event 0xc1a7df69ed0404441720a3eb710f322985ecbbe4f31934587bee703efe04427a.
//
// Solidity: event NewOracleRequest(string requestType, address callbackAddress, bytes4 callbackMethodID, uint256 requestID)
func (_Smartcontracts *SmartcontractsFilterer) ParseNewOracleRequest(log types.Log) (*SmartcontractsNewOracleRequest, error) {
	event := new(SmartcontractsNewOracleRequest)
	if err := _Smartcontracts.contract.UnpackLog(event, "NewOracleRequest", log); err != nil {
		return nil, err
	}
	return event, nil
}
