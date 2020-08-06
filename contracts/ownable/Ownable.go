// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ownable

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
const SmartcontractsABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Smartcontracts *SmartcontractsCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Smartcontracts.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Smartcontracts *SmartcontractsSession) IsOwner() (bool, error) {
	return _Smartcontracts.Contract.IsOwner(&_Smartcontracts.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_Smartcontracts *SmartcontractsCallerSession) IsOwner() (bool, error) {
	return _Smartcontracts.Contract.IsOwner(&_Smartcontracts.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Smartcontracts *SmartcontractsCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Smartcontracts.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Smartcontracts *SmartcontractsSession) Owner() (common.Address, error) {
	return _Smartcontracts.Contract.Owner(&_Smartcontracts.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Smartcontracts *SmartcontractsCallerSession) Owner() (common.Address, error) {
	return _Smartcontracts.Contract.Owner(&_Smartcontracts.CallOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Smartcontracts *SmartcontractsTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Smartcontracts.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Smartcontracts *SmartcontractsSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Smartcontracts.Contract.TransferOwnership(&_Smartcontracts.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Smartcontracts *SmartcontractsTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Smartcontracts.Contract.TransferOwnership(&_Smartcontracts.TransactOpts, newOwner)
}

// SmartcontractsOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Smartcontracts contract.
type SmartcontractsOwnershipTransferredIterator struct {
	Event *SmartcontractsOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SmartcontractsOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SmartcontractsOwnershipTransferred)
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
		it.Event = new(SmartcontractsOwnershipTransferred)
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
func (it *SmartcontractsOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SmartcontractsOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SmartcontractsOwnershipTransferred represents a OwnershipTransferred event raised by the Smartcontracts contract.
type SmartcontractsOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Smartcontracts *SmartcontractsFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SmartcontractsOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Smartcontracts.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SmartcontractsOwnershipTransferredIterator{contract: _Smartcontracts.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Smartcontracts *SmartcontractsFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SmartcontractsOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Smartcontracts.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SmartcontractsOwnershipTransferred)
				if err := _Smartcontracts.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Smartcontracts *SmartcontractsFilterer) ParseOwnershipTransferred(log types.Log) (*SmartcontractsOwnershipTransferred, error) {
	event := new(SmartcontractsOwnershipTransferred)
	if err := _Smartcontracts.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}
