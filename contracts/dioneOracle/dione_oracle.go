// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dioneOracle

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

// DioneOracleABI is the input ABI used to generate the binding from.
const DioneOracleABI = "[{\"inputs\":[{\"internalType\":\"contractIDioneStaking\",\"name\":\"_dioneStaking\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reqID\",\"type\":\"uint256\"}],\"name\":\"CancelOracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"originChain\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"requestParams\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reqID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deadline\",\"type\":\"uint256\"}],\"name\":\"NewOracleRequest\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reqID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"SubmittedOracleRequest\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"MAXIMUM_DELAY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_reqID\",\"type\":\"uint256\"}],\"name\":\"cancelOracleRequest\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dioneStaking\",\"outputs\":[{\"internalType\":\"contractIDioneStaking\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minimumDisputeFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"_originChain\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"_requestType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_requestParams\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_callbackAddress\",\"type\":\"address\"},{\"internalType\":\"bytes4\",\"name\":\"_callbackMethodID\",\"type\":\"bytes4\"}],\"name\":\"requestOracles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_reqID\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"submitOracleRequest\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// DioneOracle is an auto generated Go binding around an Ethereum contract.
type DioneOracle struct {
	DioneOracleCaller     // Read-only binding to the contract
	DioneOracleTransactor // Write-only binding to the contract
	DioneOracleFilterer   // Log filterer for contract events
}

// DioneOracleCaller is an auto generated read-only Go binding around an Ethereum contract.
type DioneOracleCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneOracleTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DioneOracleTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneOracleFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DioneOracleFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneOracleSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DioneOracleSession struct {
	Contract     *DioneOracle      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DioneOracleCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DioneOracleCallerSession struct {
	Contract *DioneOracleCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DioneOracleTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DioneOracleTransactorSession struct {
	Contract     *DioneOracleTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DioneOracleRaw is an auto generated low-level Go binding around an Ethereum contract.
type DioneOracleRaw struct {
	Contract *DioneOracle // Generic contract binding to access the raw methods on
}

// DioneOracleCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DioneOracleCallerRaw struct {
	Contract *DioneOracleCaller // Generic read-only contract binding to access the raw methods on
}

// DioneOracleTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DioneOracleTransactorRaw struct {
	Contract *DioneOracleTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDioneOracle creates a new instance of DioneOracle, bound to a specific deployed contract.
func NewDioneOracle(address common.Address, backend bind.ContractBackend) (*DioneOracle, error) {
	contract, err := bindDioneOracle(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DioneOracle{DioneOracleCaller: DioneOracleCaller{contract: contract}, DioneOracleTransactor: DioneOracleTransactor{contract: contract}, DioneOracleFilterer: DioneOracleFilterer{contract: contract}}, nil
}

// NewDioneOracleCaller creates a new read-only instance of DioneOracle, bound to a specific deployed contract.
func NewDioneOracleCaller(address common.Address, caller bind.ContractCaller) (*DioneOracleCaller, error) {
	contract, err := bindDioneOracle(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DioneOracleCaller{contract: contract}, nil
}

// NewDioneOracleTransactor creates a new write-only instance of DioneOracle, bound to a specific deployed contract.
func NewDioneOracleTransactor(address common.Address, transactor bind.ContractTransactor) (*DioneOracleTransactor, error) {
	contract, err := bindDioneOracle(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DioneOracleTransactor{contract: contract}, nil
}

// NewDioneOracleFilterer creates a new log filterer instance of DioneOracle, bound to a specific deployed contract.
func NewDioneOracleFilterer(address common.Address, filterer bind.ContractFilterer) (*DioneOracleFilterer, error) {
	contract, err := bindDioneOracle(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DioneOracleFilterer{contract: contract}, nil
}

// bindDioneOracle binds a generic wrapper to an already deployed contract.
func bindDioneOracle(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DioneOracleABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DioneOracle *DioneOracleRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DioneOracle.Contract.DioneOracleCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DioneOracle *DioneOracleRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DioneOracle.Contract.DioneOracleTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DioneOracle *DioneOracleRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DioneOracle.Contract.DioneOracleTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DioneOracle *DioneOracleCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DioneOracle.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DioneOracle *DioneOracleTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DioneOracle.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DioneOracle *DioneOracleTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DioneOracle.Contract.contract.Transact(opts, method, params...)
}

// MAXIMUMDELAY is a free data retrieval call binding the contract method 0x7d645fab.
//
// Solidity: function MAXIMUM_DELAY() view returns(uint256)
func (_DioneOracle *DioneOracleCaller) MAXIMUMDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DioneOracle.contract.Call(opts, &out, "MAXIMUM_DELAY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXIMUMDELAY is a free data retrieval call binding the contract method 0x7d645fab.
//
// Solidity: function MAXIMUM_DELAY() view returns(uint256)
func (_DioneOracle *DioneOracleSession) MAXIMUMDELAY() (*big.Int, error) {
	return _DioneOracle.Contract.MAXIMUMDELAY(&_DioneOracle.CallOpts)
}

// MAXIMUMDELAY is a free data retrieval call binding the contract method 0x7d645fab.
//
// Solidity: function MAXIMUM_DELAY() view returns(uint256)
func (_DioneOracle *DioneOracleCallerSession) MAXIMUMDELAY() (*big.Int, error) {
	return _DioneOracle.Contract.MAXIMUMDELAY(&_DioneOracle.CallOpts)
}

// DioneStaking is a free data retrieval call binding the contract method 0xe7013ddd.
//
// Solidity: function dioneStaking() view returns(address)
func (_DioneOracle *DioneOracleCaller) DioneStaking(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DioneOracle.contract.Call(opts, &out, "dioneStaking")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DioneStaking is a free data retrieval call binding the contract method 0xe7013ddd.
//
// Solidity: function dioneStaking() view returns(address)
func (_DioneOracle *DioneOracleSession) DioneStaking() (common.Address, error) {
	return _DioneOracle.Contract.DioneStaking(&_DioneOracle.CallOpts)
}

// DioneStaking is a free data retrieval call binding the contract method 0xe7013ddd.
//
// Solidity: function dioneStaking() view returns(address)
func (_DioneOracle *DioneOracleCallerSession) DioneStaking() (common.Address, error) {
	return _DioneOracle.Contract.DioneStaking(&_DioneOracle.CallOpts)
}

// MinimumDisputeFee is a free data retrieval call binding the contract method 0x4909f765.
//
// Solidity: function minimumDisputeFee() view returns(uint256)
func (_DioneOracle *DioneOracleCaller) MinimumDisputeFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DioneOracle.contract.Call(opts, &out, "minimumDisputeFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MinimumDisputeFee is a free data retrieval call binding the contract method 0x4909f765.
//
// Solidity: function minimumDisputeFee() view returns(uint256)
func (_DioneOracle *DioneOracleSession) MinimumDisputeFee() (*big.Int, error) {
	return _DioneOracle.Contract.MinimumDisputeFee(&_DioneOracle.CallOpts)
}

// MinimumDisputeFee is a free data retrieval call binding the contract method 0x4909f765.
//
// Solidity: function minimumDisputeFee() view returns(uint256)
func (_DioneOracle *DioneOracleCallerSession) MinimumDisputeFee() (*big.Int, error) {
	return _DioneOracle.Contract.MinimumDisputeFee(&_DioneOracle.CallOpts)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x1d2a198a.
//
// Solidity: function cancelOracleRequest(uint256 _reqID) returns()
func (_DioneOracle *DioneOracleTransactor) CancelOracleRequest(opts *bind.TransactOpts, _reqID *big.Int) (*types.Transaction, error) {
	return _DioneOracle.contract.Transact(opts, "cancelOracleRequest", _reqID)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x1d2a198a.
//
// Solidity: function cancelOracleRequest(uint256 _reqID) returns()
func (_DioneOracle *DioneOracleSession) CancelOracleRequest(_reqID *big.Int) (*types.Transaction, error) {
	return _DioneOracle.Contract.CancelOracleRequest(&_DioneOracle.TransactOpts, _reqID)
}

// CancelOracleRequest is a paid mutator transaction binding the contract method 0x1d2a198a.
//
// Solidity: function cancelOracleRequest(uint256 _reqID) returns()
func (_DioneOracle *DioneOracleTransactorSession) CancelOracleRequest(_reqID *big.Int) (*types.Transaction, error) {
	return _DioneOracle.Contract.CancelOracleRequest(&_DioneOracle.TransactOpts, _reqID)
}

// RequestOracles is a paid mutator transaction binding the contract method 0xe7c3712a.
//
// Solidity: function requestOracles(uint8 _originChain, string _requestType, string _requestParams, address _callbackAddress, bytes4 _callbackMethodID) returns(uint256)
func (_DioneOracle *DioneOracleTransactor) RequestOracles(opts *bind.TransactOpts, _originChain uint8, _requestType string, _requestParams string, _callbackAddress common.Address, _callbackMethodID [4]byte) (*types.Transaction, error) {
	return _DioneOracle.contract.Transact(opts, "requestOracles", _originChain, _requestType, _requestParams, _callbackAddress, _callbackMethodID)
}

// RequestOracles is a paid mutator transaction binding the contract method 0xe7c3712a.
//
// Solidity: function requestOracles(uint8 _originChain, string _requestType, string _requestParams, address _callbackAddress, bytes4 _callbackMethodID) returns(uint256)
func (_DioneOracle *DioneOracleSession) RequestOracles(_originChain uint8, _requestType string, _requestParams string, _callbackAddress common.Address, _callbackMethodID [4]byte) (*types.Transaction, error) {
	return _DioneOracle.Contract.RequestOracles(&_DioneOracle.TransactOpts, _originChain, _requestType, _requestParams, _callbackAddress, _callbackMethodID)
}

// RequestOracles is a paid mutator transaction binding the contract method 0xe7c3712a.
//
// Solidity: function requestOracles(uint8 _originChain, string _requestType, string _requestParams, address _callbackAddress, bytes4 _callbackMethodID) returns(uint256)
func (_DioneOracle *DioneOracleTransactorSession) RequestOracles(_originChain uint8, _requestType string, _requestParams string, _callbackAddress common.Address, _callbackMethodID [4]byte) (*types.Transaction, error) {
	return _DioneOracle.Contract.RequestOracles(&_DioneOracle.TransactOpts, _originChain, _requestType, _requestParams, _callbackAddress, _callbackMethodID)
}

// SubmitOracleRequest is a paid mutator transaction binding the contract method 0xcbed450e.
//
// Solidity: function submitOracleRequest(uint256 _reqID, bytes _data) returns(bool)
func (_DioneOracle *DioneOracleTransactor) SubmitOracleRequest(opts *bind.TransactOpts, _reqID *big.Int, _data []byte) (*types.Transaction, error) {
	return _DioneOracle.contract.Transact(opts, "submitOracleRequest", _reqID, _data)
}

// SubmitOracleRequest is a paid mutator transaction binding the contract method 0xcbed450e.
//
// Solidity: function submitOracleRequest(uint256 _reqID, bytes _data) returns(bool)
func (_DioneOracle *DioneOracleSession) SubmitOracleRequest(_reqID *big.Int, _data []byte) (*types.Transaction, error) {
	return _DioneOracle.Contract.SubmitOracleRequest(&_DioneOracle.TransactOpts, _reqID, _data)
}

// SubmitOracleRequest is a paid mutator transaction binding the contract method 0xcbed450e.
//
// Solidity: function submitOracleRequest(uint256 _reqID, bytes _data) returns(bool)
func (_DioneOracle *DioneOracleTransactorSession) SubmitOracleRequest(_reqID *big.Int, _data []byte) (*types.Transaction, error) {
	return _DioneOracle.Contract.SubmitOracleRequest(&_DioneOracle.TransactOpts, _reqID, _data)
}

// DioneOracleCancelOracleRequestIterator is returned from FilterCancelOracleRequest and is used to iterate over the raw logs and unpacked data for CancelOracleRequest events raised by the DioneOracle contract.
type DioneOracleCancelOracleRequestIterator struct {
	Event *DioneOracleCancelOracleRequest // Event containing the contract specifics and raw log

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
func (it *DioneOracleCancelOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneOracleCancelOracleRequest)
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
		it.Event = new(DioneOracleCancelOracleRequest)
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
func (it *DioneOracleCancelOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneOracleCancelOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneOracleCancelOracleRequest represents a CancelOracleRequest event raised by the DioneOracle contract.
type DioneOracleCancelOracleRequest struct {
	ReqID *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterCancelOracleRequest is a free log retrieval operation binding the contract event 0xb9422b5075524babd23853e336b5ef3a05516a8088f4c8b809f949642f6aa54e.
//
// Solidity: event CancelOracleRequest(uint256 reqID)
func (_DioneOracle *DioneOracleFilterer) FilterCancelOracleRequest(opts *bind.FilterOpts) (*DioneOracleCancelOracleRequestIterator, error) {

	logs, sub, err := _DioneOracle.contract.FilterLogs(opts, "CancelOracleRequest")
	if err != nil {
		return nil, err
	}
	return &DioneOracleCancelOracleRequestIterator{contract: _DioneOracle.contract, event: "CancelOracleRequest", logs: logs, sub: sub}, nil
}

// WatchCancelOracleRequest is a free log subscription operation binding the contract event 0xb9422b5075524babd23853e336b5ef3a05516a8088f4c8b809f949642f6aa54e.
//
// Solidity: event CancelOracleRequest(uint256 reqID)
func (_DioneOracle *DioneOracleFilterer) WatchCancelOracleRequest(opts *bind.WatchOpts, sink chan<- *DioneOracleCancelOracleRequest) (event.Subscription, error) {

	logs, sub, err := _DioneOracle.contract.WatchLogs(opts, "CancelOracleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneOracleCancelOracleRequest)
				if err := _DioneOracle.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
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

// ParseCancelOracleRequest is a log parse operation binding the contract event 0xb9422b5075524babd23853e336b5ef3a05516a8088f4c8b809f949642f6aa54e.
//
// Solidity: event CancelOracleRequest(uint256 reqID)
func (_DioneOracle *DioneOracleFilterer) ParseCancelOracleRequest(log types.Log) (*DioneOracleCancelOracleRequest, error) {
	event := new(DioneOracleCancelOracleRequest)
	if err := _DioneOracle.contract.UnpackLog(event, "CancelOracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DioneOracleNewOracleRequestIterator is returned from FilterNewOracleRequest and is used to iterate over the raw logs and unpacked data for NewOracleRequest events raised by the DioneOracle contract.
type DioneOracleNewOracleRequestIterator struct {
	Event *DioneOracleNewOracleRequest // Event containing the contract specifics and raw log

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
func (it *DioneOracleNewOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneOracleNewOracleRequest)
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
		it.Event = new(DioneOracleNewOracleRequest)
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
func (it *DioneOracleNewOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneOracleNewOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneOracleNewOracleRequest represents a NewOracleRequest event raised by the DioneOracle contract.
type DioneOracleNewOracleRequest struct {
	OriginChain   uint8
	RequestType   string
	RequestParams string
	ReqID         *big.Int
	Deadline      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterNewOracleRequest is a free log retrieval operation binding the contract event 0xf2736922e09befc7bb0c19052cc9b530a5d9a696b3db44a29814ac8a4335aa08.
//
// Solidity: event NewOracleRequest(uint8 originChain, string requestType, string requestParams, uint256 reqID, uint256 deadline)
func (_DioneOracle *DioneOracleFilterer) FilterNewOracleRequest(opts *bind.FilterOpts) (*DioneOracleNewOracleRequestIterator, error) {

	logs, sub, err := _DioneOracle.contract.FilterLogs(opts, "NewOracleRequest")
	if err != nil {
		return nil, err
	}
	return &DioneOracleNewOracleRequestIterator{contract: _DioneOracle.contract, event: "NewOracleRequest", logs: logs, sub: sub}, nil
}

// WatchNewOracleRequest is a free log subscription operation binding the contract event 0xf2736922e09befc7bb0c19052cc9b530a5d9a696b3db44a29814ac8a4335aa08.
//
// Solidity: event NewOracleRequest(uint8 originChain, string requestType, string requestParams, uint256 reqID, uint256 deadline)
func (_DioneOracle *DioneOracleFilterer) WatchNewOracleRequest(opts *bind.WatchOpts, sink chan<- *DioneOracleNewOracleRequest) (event.Subscription, error) {

	logs, sub, err := _DioneOracle.contract.WatchLogs(opts, "NewOracleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneOracleNewOracleRequest)
				if err := _DioneOracle.contract.UnpackLog(event, "NewOracleRequest", log); err != nil {
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

// ParseNewOracleRequest is a log parse operation binding the contract event 0xf2736922e09befc7bb0c19052cc9b530a5d9a696b3db44a29814ac8a4335aa08.
//
// Solidity: event NewOracleRequest(uint8 originChain, string requestType, string requestParams, uint256 reqID, uint256 deadline)
func (_DioneOracle *DioneOracleFilterer) ParseNewOracleRequest(log types.Log) (*DioneOracleNewOracleRequest, error) {
	event := new(DioneOracleNewOracleRequest)
	if err := _DioneOracle.contract.UnpackLog(event, "NewOracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DioneOracleSubmittedOracleRequestIterator is returned from FilterSubmittedOracleRequest and is used to iterate over the raw logs and unpacked data for SubmittedOracleRequest events raised by the DioneOracle contract.
type DioneOracleSubmittedOracleRequestIterator struct {
	Event *DioneOracleSubmittedOracleRequest // Event containing the contract specifics and raw log

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
func (it *DioneOracleSubmittedOracleRequestIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneOracleSubmittedOracleRequest)
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
		it.Event = new(DioneOracleSubmittedOracleRequest)
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
func (it *DioneOracleSubmittedOracleRequestIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneOracleSubmittedOracleRequestIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneOracleSubmittedOracleRequest represents a SubmittedOracleRequest event raised by the DioneOracle contract.
type DioneOracleSubmittedOracleRequest struct {
	ReqID *big.Int
	Data  []byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterSubmittedOracleRequest is a free log retrieval operation binding the contract event 0x30c6663e7dc6e4e92bde7b263ac094ed9141eba4437261c482c14e41d855040f.
//
// Solidity: event SubmittedOracleRequest(uint256 reqID, bytes data)
func (_DioneOracle *DioneOracleFilterer) FilterSubmittedOracleRequest(opts *bind.FilterOpts) (*DioneOracleSubmittedOracleRequestIterator, error) {

	logs, sub, err := _DioneOracle.contract.FilterLogs(opts, "SubmittedOracleRequest")
	if err != nil {
		return nil, err
	}
	return &DioneOracleSubmittedOracleRequestIterator{contract: _DioneOracle.contract, event: "SubmittedOracleRequest", logs: logs, sub: sub}, nil
}

// WatchSubmittedOracleRequest is a free log subscription operation binding the contract event 0x30c6663e7dc6e4e92bde7b263ac094ed9141eba4437261c482c14e41d855040f.
//
// Solidity: event SubmittedOracleRequest(uint256 reqID, bytes data)
func (_DioneOracle *DioneOracleFilterer) WatchSubmittedOracleRequest(opts *bind.WatchOpts, sink chan<- *DioneOracleSubmittedOracleRequest) (event.Subscription, error) {

	logs, sub, err := _DioneOracle.contract.WatchLogs(opts, "SubmittedOracleRequest")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneOracleSubmittedOracleRequest)
				if err := _DioneOracle.contract.UnpackLog(event, "SubmittedOracleRequest", log); err != nil {
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

// ParseSubmittedOracleRequest is a log parse operation binding the contract event 0x30c6663e7dc6e4e92bde7b263ac094ed9141eba4437261c482c14e41d855040f.
//
// Solidity: event SubmittedOracleRequest(uint256 reqID, bytes data)
func (_DioneOracle *DioneOracleFilterer) ParseSubmittedOracleRequest(log types.Log) (*DioneOracleSubmittedOracleRequest, error) {
	event := new(DioneOracleSubmittedOracleRequest)
	if err := _DioneOracle.contract.UnpackLog(event, "SubmittedOracleRequest", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
