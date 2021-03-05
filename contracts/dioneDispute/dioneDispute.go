// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dioneDispute

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

// DioneDisputeABI is the input ABI used to generate the binding from.
const DioneDisputeABI = "[{\"inputs\":[{\"internalType\":\"contractIDioneStaking\",\"name\":\"_dioneStaking\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dhash\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"status\",\"type\":\"bool\"}],\"name\":\"DisputeFinished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dhash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"miner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"disputeInitiator\",\"type\":\"address\"}],\"name\":\"NewDispute\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dhash\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"votedMiner\",\"type\":\"address\"}],\"name\":\"NewVote\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"dioneStaking\",\"outputs\":[{\"internalType\":\"contractIDioneStaking\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"miner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"beginDispute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dhash\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"voteStatus\",\"type\":\"bool\"}],\"name\":\"vote\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"dhash\",\"type\":\"bytes32\"}],\"name\":\"finishDispute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}],"

// DioneDispute is an auto generated Go binding around an Ethereum contract.
type DioneDispute struct {
	DioneDisputeCaller     // Read-only binding to the contract
	DioneDisputeTransactor // Write-only binding to the contract
	DioneDisputeFilterer   // Log filterer for contract events
}

// DioneDisputeCaller is an auto generated read-only Go binding around an Ethereum contract.
type DioneDisputeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneDisputeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DioneDisputeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneDisputeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DioneDisputeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneDisputeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DioneDisputeSession struct {
	Contract     *DioneDispute     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DioneDisputeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DioneDisputeCallerSession struct {
	Contract *DioneDisputeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// DioneDisputeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DioneDisputeTransactorSession struct {
	Contract     *DioneDisputeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// DioneDisputeRaw is an auto generated low-level Go binding around an Ethereum contract.
type DioneDisputeRaw struct {
	Contract *DioneDispute // Generic contract binding to access the raw methods on
}

// DioneDisputeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DioneDisputeCallerRaw struct {
	Contract *DioneDisputeCaller // Generic read-only contract binding to access the raw methods on
}

// DioneDisputeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DioneDisputeTransactorRaw struct {
	Contract *DioneDisputeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDioneDispute creates a new instance of DioneDispute, bound to a specific deployed contract.
func NewDioneDispute(address common.Address, backend bind.ContractBackend) (*DioneDispute, error) {
	contract, err := bindDioneDispute(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DioneDispute{DioneDisputeCaller: DioneDisputeCaller{contract: contract}, DioneDisputeTransactor: DioneDisputeTransactor{contract: contract}, DioneDisputeFilterer: DioneDisputeFilterer{contract: contract}}, nil
}

// NewDioneDisputeCaller creates a new read-only instance of DioneDispute, bound to a specific deployed contract.
func NewDioneDisputeCaller(address common.Address, caller bind.ContractCaller) (*DioneDisputeCaller, error) {
	contract, err := bindDioneDispute(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DioneDisputeCaller{contract: contract}, nil
}

// NewDioneDisputeTransactor creates a new write-only instance of DioneDispute, bound to a specific deployed contract.
func NewDioneDisputeTransactor(address common.Address, transactor bind.ContractTransactor) (*DioneDisputeTransactor, error) {
	contract, err := bindDioneDispute(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DioneDisputeTransactor{contract: contract}, nil
}

// NewDioneDisputeFilterer creates a new log filterer instance of DioneDispute, bound to a specific deployed contract.
func NewDioneDisputeFilterer(address common.Address, filterer bind.ContractFilterer) (*DioneDisputeFilterer, error) {
	contract, err := bindDioneDispute(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DioneDisputeFilterer{contract: contract}, nil
}

// bindDioneDispute binds a generic wrapper to an already deployed contract.
func bindDioneDispute(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DioneDisputeABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DioneDispute *DioneDisputeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DioneDispute.Contract.DioneDisputeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DioneDispute *DioneDisputeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DioneDispute.Contract.DioneDisputeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DioneDispute *DioneDisputeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DioneDispute.Contract.DioneDisputeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DioneDispute *DioneDisputeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DioneDispute.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DioneDispute *DioneDisputeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DioneDispute.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DioneDispute *DioneDisputeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DioneDispute.Contract.contract.Transact(opts, method, params...)
}

// DioneStaking is a free data retrieval call binding the contract method 0xe7013ddd.
//
// Solidity: function dioneStaking() view returns(address)
func (_DioneDispute *DioneDisputeCaller) DioneStaking(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DioneDispute.contract.Call(opts, &out, "dioneStaking")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DioneStaking is a free data retrieval call binding the contract method 0xe7013ddd.
//
// Solidity: function dioneStaking() view returns(address)
func (_DioneDispute *DioneDisputeSession) DioneStaking() (common.Address, error) {
	return _DioneDispute.Contract.DioneStaking(&_DioneDispute.CallOpts)
}

// DioneStaking is a free data retrieval call binding the contract method 0xe7013ddd.
//
// Solidity: function dioneStaking() view returns(address)
func (_DioneDispute *DioneDisputeCallerSession) DioneStaking() (common.Address, error) {
	return _DioneDispute.Contract.DioneStaking(&_DioneDispute.CallOpts)
}

// BeginDispute is a paid mutator transaction binding the contract method 0x0b5e2057.
//
// Solidity: function beginDispute(address miner, uint256 requestID) returns()
func (_DioneDispute *DioneDisputeTransactor) BeginDispute(opts *bind.TransactOpts, miner common.Address, requestID *big.Int) (*types.Transaction, error) {
	return _DioneDispute.contract.Transact(opts, "beginDispute", miner, requestID)
}

// BeginDispute is a paid mutator transaction binding the contract method 0x0b5e2057.
//
// Solidity: function beginDispute(address miner, uint256 requestID) returns()
func (_DioneDispute *DioneDisputeSession) BeginDispute(miner common.Address, requestID *big.Int) (*types.Transaction, error) {
	return _DioneDispute.Contract.BeginDispute(&_DioneDispute.TransactOpts, miner, requestID)
}

// BeginDispute is a paid mutator transaction binding the contract method 0x0b5e2057.
//
// Solidity: function beginDispute(address miner, uint256 requestID) returns()
func (_DioneDispute *DioneDisputeTransactorSession) BeginDispute(miner common.Address, requestID *big.Int) (*types.Transaction, error) {
	return _DioneDispute.Contract.BeginDispute(&_DioneDispute.TransactOpts, miner, requestID)
}

// FinishDispute is a paid mutator transaction binding the contract method 0xa597d7c5.
//
// Solidity: function finishDispute(bytes32 dhash) returns()
func (_DioneDispute *DioneDisputeTransactor) FinishDispute(opts *bind.TransactOpts, dhash [32]byte) (*types.Transaction, error) {
	return _DioneDispute.contract.Transact(opts, "finishDispute", dhash)
}

// FinishDispute is a paid mutator transaction binding the contract method 0xa597d7c5.
//
// Solidity: function finishDispute(bytes32 dhash) returns()
func (_DioneDispute *DioneDisputeSession) FinishDispute(dhash [32]byte) (*types.Transaction, error) {
	return _DioneDispute.Contract.FinishDispute(&_DioneDispute.TransactOpts, dhash)
}

// FinishDispute is a paid mutator transaction binding the contract method 0xa597d7c5.
//
// Solidity: function finishDispute(bytes32 dhash) returns()
func (_DioneDispute *DioneDisputeTransactorSession) FinishDispute(dhash [32]byte) (*types.Transaction, error) {
	return _DioneDispute.Contract.FinishDispute(&_DioneDispute.TransactOpts, dhash)
}

// Vote is a paid mutator transaction binding the contract method 0x9f2ce678.
//
// Solidity: function vote(bytes32 dhash, bool voteStatus) returns()
func (_DioneDispute *DioneDisputeTransactor) Vote(opts *bind.TransactOpts, dhash [32]byte, voteStatus bool) (*types.Transaction, error) {
	return _DioneDispute.contract.Transact(opts, "vote", dhash, voteStatus)
}

// Vote is a paid mutator transaction binding the contract method 0x9f2ce678.
//
// Solidity: function vote(bytes32 dhash, bool voteStatus) returns()
func (_DioneDispute *DioneDisputeSession) Vote(dhash [32]byte, voteStatus bool) (*types.Transaction, error) {
	return _DioneDispute.Contract.Vote(&_DioneDispute.TransactOpts, dhash, voteStatus)
}

// Vote is a paid mutator transaction binding the contract method 0x9f2ce678.
//
// Solidity: function vote(bytes32 dhash, bool voteStatus) returns()
func (_DioneDispute *DioneDisputeTransactorSession) Vote(dhash [32]byte, voteStatus bool) (*types.Transaction, error) {
	return _DioneDispute.Contract.Vote(&_DioneDispute.TransactOpts, dhash, voteStatus)
}

// DioneDisputeDisputeFinishedIterator is returned from FilterDisputeFinished and is used to iterate over the raw logs and unpacked data for DisputeFinished events raised by the DioneDispute contract.
type DioneDisputeDisputeFinishedIterator struct {
	Event *DioneDisputeDisputeFinished // Event containing the contract specifics and raw log

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
func (it *DioneDisputeDisputeFinishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneDisputeDisputeFinished)
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
		it.Event = new(DioneDisputeDisputeFinished)
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
func (it *DioneDisputeDisputeFinishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneDisputeDisputeFinishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneDisputeDisputeFinished represents a DisputeFinished event raised by the DioneDispute contract.
type DioneDisputeDisputeFinished struct {
	Dhash  [32]byte
	Status bool
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDisputeFinished is a free log retrieval operation binding the contract event 0x4d9456641532b60dd6621cc1147a253cb157757d8df5b07c92aab0e82d4c3359.
//
// Solidity: event DisputeFinished(bytes32 dhash, bool status)
func (_DioneDispute *DioneDisputeFilterer) FilterDisputeFinished(opts *bind.FilterOpts) (*DioneDisputeDisputeFinishedIterator, error) {

	logs, sub, err := _DioneDispute.contract.FilterLogs(opts, "DisputeFinished")
	if err != nil {
		return nil, err
	}
	return &DioneDisputeDisputeFinishedIterator{contract: _DioneDispute.contract, event: "DisputeFinished", logs: logs, sub: sub}, nil
}

// WatchDisputeFinished is a free log subscription operation binding the contract event 0x4d9456641532b60dd6621cc1147a253cb157757d8df5b07c92aab0e82d4c3359.
//
// Solidity: event DisputeFinished(bytes32 dhash, bool status)
func (_DioneDispute *DioneDisputeFilterer) WatchDisputeFinished(opts *bind.WatchOpts, sink chan<- *DioneDisputeDisputeFinished) (event.Subscription, error) {

	logs, sub, err := _DioneDispute.contract.WatchLogs(opts, "DisputeFinished")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneDisputeDisputeFinished)
				if err := _DioneDispute.contract.UnpackLog(event, "DisputeFinished", log); err != nil {
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

// ParseDisputeFinished is a log parse operation binding the contract event 0x4d9456641532b60dd6621cc1147a253cb157757d8df5b07c92aab0e82d4c3359.
//
// Solidity: event DisputeFinished(bytes32 dhash, bool status)
func (_DioneDispute *DioneDisputeFilterer) ParseDisputeFinished(log types.Log) (*DioneDisputeDisputeFinished, error) {
	event := new(DioneDisputeDisputeFinished)
	if err := _DioneDispute.contract.UnpackLog(event, "DisputeFinished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DioneDisputeNewDisputeIterator is returned from FilterNewDispute and is used to iterate over the raw logs and unpacked data for NewDispute events raised by the DioneDispute contract.
type DioneDisputeNewDisputeIterator struct {
	Event *DioneDisputeNewDispute // Event containing the contract specifics and raw log

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
func (it *DioneDisputeNewDisputeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneDisputeNewDispute)
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
		it.Event = new(DioneDisputeNewDispute)
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
func (it *DioneDisputeNewDisputeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneDisputeNewDisputeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneDisputeNewDispute represents a NewDispute event raised by the DioneDispute contract.
type DioneDisputeNewDispute struct {
	Dhash            [32]byte
	Miner            common.Address
	DisputeInitiator common.Address
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterNewDispute is a free log retrieval operation binding the contract event 0xcd3d2ebd72bbbf7e5659c6790b60f42174aafe2741b7ba3807007455faae522e.
//
// Solidity: event NewDispute(bytes32 dhash, address indexed miner, address indexed disputeInitiator)
func (_DioneDispute *DioneDisputeFilterer) FilterNewDispute(opts *bind.FilterOpts, miner []common.Address, disputeInitiator []common.Address) (*DioneDisputeNewDisputeIterator, error) {

	var minerRule []interface{}
	for _, minerItem := range miner {
		minerRule = append(minerRule, minerItem)
	}
	var disputeInitiatorRule []interface{}
	for _, disputeInitiatorItem := range disputeInitiator {
		disputeInitiatorRule = append(disputeInitiatorRule, disputeInitiatorItem)
	}

	logs, sub, err := _DioneDispute.contract.FilterLogs(opts, "NewDispute", minerRule, disputeInitiatorRule)
	if err != nil {
		return nil, err
	}
	return &DioneDisputeNewDisputeIterator{contract: _DioneDispute.contract, event: "NewDispute", logs: logs, sub: sub}, nil
}

// WatchNewDispute is a free log subscription operation binding the contract event 0xcd3d2ebd72bbbf7e5659c6790b60f42174aafe2741b7ba3807007455faae522e.
//
// Solidity: event NewDispute(bytes32 dhash, address indexed miner, address indexed disputeInitiator)
func (_DioneDispute *DioneDisputeFilterer) WatchNewDispute(opts *bind.WatchOpts, sink chan<- *DioneDisputeNewDispute, miner []common.Address, disputeInitiator []common.Address) (event.Subscription, error) {

	var minerRule []interface{}
	for _, minerItem := range miner {
		minerRule = append(minerRule, minerItem)
	}
	var disputeInitiatorRule []interface{}
	for _, disputeInitiatorItem := range disputeInitiator {
		disputeInitiatorRule = append(disputeInitiatorRule, disputeInitiatorItem)
	}

	logs, sub, err := _DioneDispute.contract.WatchLogs(opts, "NewDispute", minerRule, disputeInitiatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneDisputeNewDispute)
				if err := _DioneDispute.contract.UnpackLog(event, "NewDispute", log); err != nil {
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

// ParseNewDispute is a log parse operation binding the contract event 0xcd3d2ebd72bbbf7e5659c6790b60f42174aafe2741b7ba3807007455faae522e.
//
// Solidity: event NewDispute(bytes32 dhash, address indexed miner, address indexed disputeInitiator)
func (_DioneDispute *DioneDisputeFilterer) ParseNewDispute(log types.Log) (*DioneDisputeNewDispute, error) {
	event := new(DioneDisputeNewDispute)
	if err := _DioneDispute.contract.UnpackLog(event, "NewDispute", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DioneDisputeNewVoteIterator is returned from FilterNewVote and is used to iterate over the raw logs and unpacked data for NewVote events raised by the DioneDispute contract.
type DioneDisputeNewVoteIterator struct {
	Event *DioneDisputeNewVote // Event containing the contract specifics and raw log

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
func (it *DioneDisputeNewVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneDisputeNewVote)
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
		it.Event = new(DioneDisputeNewVote)
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
func (it *DioneDisputeNewVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneDisputeNewVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneDisputeNewVote represents a NewVote event raised by the DioneDispute contract.
type DioneDisputeNewVote struct {
	Dhash      [32]byte
	VotedMiner common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterNewVote is a free log retrieval operation binding the contract event 0xc0f8946b28ac7993b4dd24b05482b70709d3394c92247236929d86a50b24d0ab.
//
// Solidity: event NewVote(bytes32 dhash, address indexed votedMiner)
func (_DioneDispute *DioneDisputeFilterer) FilterNewVote(opts *bind.FilterOpts, votedMiner []common.Address) (*DioneDisputeNewVoteIterator, error) {

	var votedMinerRule []interface{}
	for _, votedMinerItem := range votedMiner {
		votedMinerRule = append(votedMinerRule, votedMinerItem)
	}

	logs, sub, err := _DioneDispute.contract.FilterLogs(opts, "NewVote", votedMinerRule)
	if err != nil {
		return nil, err
	}
	return &DioneDisputeNewVoteIterator{contract: _DioneDispute.contract, event: "NewVote", logs: logs, sub: sub}, nil
}

// WatchNewVote is a free log subscription operation binding the contract event 0xc0f8946b28ac7993b4dd24b05482b70709d3394c92247236929d86a50b24d0ab.
//
// Solidity: event NewVote(bytes32 dhash, address indexed votedMiner)
func (_DioneDispute *DioneDisputeFilterer) WatchNewVote(opts *bind.WatchOpts, sink chan<- *DioneDisputeNewVote, votedMiner []common.Address) (event.Subscription, error) {

	var votedMinerRule []interface{}
	for _, votedMinerItem := range votedMiner {
		votedMinerRule = append(votedMinerRule, votedMinerItem)
	}

	logs, sub, err := _DioneDispute.contract.WatchLogs(opts, "NewVote", votedMinerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneDisputeNewVote)
				if err := _DioneDispute.contract.UnpackLog(event, "NewVote", log); err != nil {
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

// ParseNewVote is a log parse operation binding the contract event 0xc0f8946b28ac7993b4dd24b05482b70709d3394c92247236929d86a50b24d0ab.
//
// Solidity: event NewVote(bytes32 dhash, address indexed votedMiner)
func (_DioneDispute *DioneDisputeFilterer) ParseNewVote(log types.Log) (*DioneDisputeNewVote, error) {
	event := new(DioneDisputeNewVote)
	if err := _DioneDispute.contract.UnpackLog(event, "NewVote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
