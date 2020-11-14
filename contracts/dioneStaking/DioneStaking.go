// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dioneStaking

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

// DioneStakingABI is the input ABI used to generate the binding from.
const DioneStakingABI = "[{\"inputs\":[{\"internalType\":\"contractDioneToken\",\"name\":\"_dione\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_aggregatorAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_minerReward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_startBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_minimumStake\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"miner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"Mine\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"miner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Stake\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"miner\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"aggregatorAddr\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dione\",\"outputs\":[{\"internalType\":\"contractDioneToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"minerInfo\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"firstStakeBlock\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"lastRewardBlock\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minerReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"minimumStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"startBlock\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_minerAddr\",\"type\":\"address\"}],\"name\":\"mine\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_minerAddr\",\"type\":\"address\"}],\"name\":\"mineAndStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_minerAddr\",\"type\":\"address\"}],\"name\":\"minerStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minerReward\",\"type\":\"uint256\"}],\"name\":\"setMinerReward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_minerAddr\",\"type\":\"address\"}],\"name\":\"isLegitMiner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_minimumStake\",\"type\":\"uint256\"}],\"name\":\"setMinimumStake\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// DioneStaking is an auto generated Go binding around an Ethereum contract.
type DioneStaking struct {
	DioneStakingCaller     // Read-only binding to the contract
	DioneStakingTransactor // Write-only binding to the contract
	DioneStakingFilterer   // Log filterer for contract events
}

// DioneStakingCaller is an auto generated read-only Go binding around an Ethereum contract.
type DioneStakingCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneStakingTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DioneStakingTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneStakingFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DioneStakingFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DioneStakingSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DioneStakingSession struct {
	Contract     *DioneStaking     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DioneStakingCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DioneStakingCallerSession struct {
	Contract *DioneStakingCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// DioneStakingTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DioneStakingTransactorSession struct {
	Contract     *DioneStakingTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// DioneStakingRaw is an auto generated low-level Go binding around an Ethereum contract.
type DioneStakingRaw struct {
	Contract *DioneStaking // Generic contract binding to access the raw methods on
}

// DioneStakingCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DioneStakingCallerRaw struct {
	Contract *DioneStakingCaller // Generic read-only contract binding to access the raw methods on
}

// DioneStakingTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DioneStakingTransactorRaw struct {
	Contract *DioneStakingTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDioneStaking creates a new instance of DioneStaking, bound to a specific deployed contract.
func NewDioneStaking(address common.Address, backend bind.ContractBackend) (*DioneStaking, error) {
	contract, err := bindDioneStaking(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DioneStaking{DioneStakingCaller: DioneStakingCaller{contract: contract}, DioneStakingTransactor: DioneStakingTransactor{contract: contract}, DioneStakingFilterer: DioneStakingFilterer{contract: contract}}, nil
}

// NewDioneStakingCaller creates a new read-only instance of DioneStaking, bound to a specific deployed contract.
func NewDioneStakingCaller(address common.Address, caller bind.ContractCaller) (*DioneStakingCaller, error) {
	contract, err := bindDioneStaking(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DioneStakingCaller{contract: contract}, nil
}

// NewDioneStakingTransactor creates a new write-only instance of DioneStaking, bound to a specific deployed contract.
func NewDioneStakingTransactor(address common.Address, transactor bind.ContractTransactor) (*DioneStakingTransactor, error) {
	contract, err := bindDioneStaking(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DioneStakingTransactor{contract: contract}, nil
}

// NewDioneStakingFilterer creates a new log filterer instance of DioneStaking, bound to a specific deployed contract.
func NewDioneStakingFilterer(address common.Address, filterer bind.ContractFilterer) (*DioneStakingFilterer, error) {
	contract, err := bindDioneStaking(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DioneStakingFilterer{contract: contract}, nil
}

// bindDioneStaking binds a generic wrapper to an already deployed contract.
func bindDioneStaking(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DioneStakingABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DioneStaking *DioneStakingRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DioneStaking.Contract.DioneStakingCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DioneStaking *DioneStakingRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DioneStaking.Contract.DioneStakingTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DioneStaking *DioneStakingRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DioneStaking.Contract.DioneStakingTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DioneStaking *DioneStakingCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DioneStaking.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DioneStaking *DioneStakingTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DioneStaking.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DioneStaking *DioneStakingTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DioneStaking.Contract.contract.Transact(opts, method, params...)
}

// AggregatorAddr is a free data retrieval call binding the contract method 0x82762600.
//
// Solidity: function aggregatorAddr() view returns(address)
func (_DioneStaking *DioneStakingCaller) AggregatorAddr(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DioneStaking.contract.Call(opts, out, "aggregatorAddr")
	return *ret0, err
}

// AggregatorAddr is a free data retrieval call binding the contract method 0x82762600.
//
// Solidity: function aggregatorAddr() view returns(address)
func (_DioneStaking *DioneStakingSession) AggregatorAddr() (common.Address, error) {
	return _DioneStaking.Contract.AggregatorAddr(&_DioneStaking.CallOpts)
}

// AggregatorAddr is a free data retrieval call binding the contract method 0x82762600.
//
// Solidity: function aggregatorAddr() view returns(address)
func (_DioneStaking *DioneStakingCallerSession) AggregatorAddr() (common.Address, error) {
	return _DioneStaking.Contract.AggregatorAddr(&_DioneStaking.CallOpts)
}

// Dione is a free data retrieval call binding the contract method 0x3425dfa6.
//
// Solidity: function dione() view returns(address)
func (_DioneStaking *DioneStakingCaller) Dione(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DioneStaking.contract.Call(opts, out, "dione")
	return *ret0, err
}

// Dione is a free data retrieval call binding the contract method 0x3425dfa6.
//
// Solidity: function dione() view returns(address)
func (_DioneStaking *DioneStakingSession) Dione() (common.Address, error) {
	return _DioneStaking.Contract.Dione(&_DioneStaking.CallOpts)
}

// Dione is a free data retrieval call binding the contract method 0x3425dfa6.
//
// Solidity: function dione() view returns(address)
func (_DioneStaking *DioneStakingCallerSession) Dione() (common.Address, error) {
	return _DioneStaking.Contract.Dione(&_DioneStaking.CallOpts)
}

// MinerInfo is a free data retrieval call binding the contract method 0x03337fd8.
//
// Solidity: function minerInfo(address ) view returns(uint256 amount, uint256 firstStakeBlock, uint256 lastRewardBlock)
func (_DioneStaking *DioneStakingCaller) MinerInfo(opts *bind.CallOpts, arg0 common.Address) (struct {
	Amount          *big.Int
	FirstStakeBlock *big.Int
	LastRewardBlock *big.Int
}, error) {
	ret := new(struct {
		Amount          *big.Int
		FirstStakeBlock *big.Int
		LastRewardBlock *big.Int
	})
	out := ret
	err := _DioneStaking.contract.Call(opts, out, "minerInfo", arg0)
	return *ret, err
}

// MinerInfo is a free data retrieval call binding the contract method 0x03337fd8.
//
// Solidity: function minerInfo(address ) view returns(uint256 amount, uint256 firstStakeBlock, uint256 lastRewardBlock)
func (_DioneStaking *DioneStakingSession) MinerInfo(arg0 common.Address) (struct {
	Amount          *big.Int
	FirstStakeBlock *big.Int
	LastRewardBlock *big.Int
}, error) {
	return _DioneStaking.Contract.MinerInfo(&_DioneStaking.CallOpts, arg0)
}

// MinerInfo is a free data retrieval call binding the contract method 0x03337fd8.
//
// Solidity: function minerInfo(address ) view returns(uint256 amount, uint256 firstStakeBlock, uint256 lastRewardBlock)
func (_DioneStaking *DioneStakingCallerSession) MinerInfo(arg0 common.Address) (struct {
	Amount          *big.Int
	FirstStakeBlock *big.Int
	LastRewardBlock *big.Int
}, error) {
	return _DioneStaking.Contract.MinerInfo(&_DioneStaking.CallOpts, arg0)
}

// MinerReward is a free data retrieval call binding the contract method 0xcbed45eb.
//
// Solidity: function minerReward() view returns(uint256)
func (_DioneStaking *DioneStakingCaller) MinerReward(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DioneStaking.contract.Call(opts, out, "minerReward")
	return *ret0, err
}

// MinerReward is a free data retrieval call binding the contract method 0xcbed45eb.
//
// Solidity: function minerReward() view returns(uint256)
func (_DioneStaking *DioneStakingSession) MinerReward() (*big.Int, error) {
	return _DioneStaking.Contract.MinerReward(&_DioneStaking.CallOpts)
}

// MinerReward is a free data retrieval call binding the contract method 0xcbed45eb.
//
// Solidity: function minerReward() view returns(uint256)
func (_DioneStaking *DioneStakingCallerSession) MinerReward() (*big.Int, error) {
	return _DioneStaking.Contract.MinerReward(&_DioneStaking.CallOpts)
}

// MinerStake is a free data retrieval call binding the contract method 0x8eaa3850.
//
// Solidity: function minerStake(address _minerAddr) view returns(uint256)
func (_DioneStaking *DioneStakingCaller) MinerStake(opts *bind.CallOpts, _minerAddr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DioneStaking.contract.Call(opts, out, "minerStake", _minerAddr)
	return *ret0, err
}

// MinerStake is a free data retrieval call binding the contract method 0x8eaa3850.
//
// Solidity: function minerStake(address _minerAddr) view returns(uint256)
func (_DioneStaking *DioneStakingSession) MinerStake(_minerAddr common.Address) (*big.Int, error) {
	return _DioneStaking.Contract.MinerStake(&_DioneStaking.CallOpts, _minerAddr)
}

// MinerStake is a free data retrieval call binding the contract method 0x8eaa3850.
//
// Solidity: function minerStake(address _minerAddr) view returns(uint256)
func (_DioneStaking *DioneStakingCallerSession) MinerStake(_minerAddr common.Address) (*big.Int, error) {
	return _DioneStaking.Contract.MinerStake(&_DioneStaking.CallOpts, _minerAddr)
}

// MinimumStake is a free data retrieval call binding the contract method 0xec5ffac2.
//
// Solidity: function minimumStake() view returns(uint256)
func (_DioneStaking *DioneStakingCaller) MinimumStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DioneStaking.contract.Call(opts, out, "minimumStake")
	return *ret0, err
}

// MinimumStake is a free data retrieval call binding the contract method 0xec5ffac2.
//
// Solidity: function minimumStake() view returns(uint256)
func (_DioneStaking *DioneStakingSession) MinimumStake() (*big.Int, error) {
	return _DioneStaking.Contract.MinimumStake(&_DioneStaking.CallOpts)
}

// MinimumStake is a free data retrieval call binding the contract method 0xec5ffac2.
//
// Solidity: function minimumStake() view returns(uint256)
func (_DioneStaking *DioneStakingCallerSession) MinimumStake() (*big.Int, error) {
	return _DioneStaking.Contract.MinimumStake(&_DioneStaking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DioneStaking *DioneStakingCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DioneStaking.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DioneStaking *DioneStakingSession) Owner() (common.Address, error) {
	return _DioneStaking.Contract.Owner(&_DioneStaking.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_DioneStaking *DioneStakingCallerSession) Owner() (common.Address, error) {
	return _DioneStaking.Contract.Owner(&_DioneStaking.CallOpts)
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() view returns(uint256)
func (_DioneStaking *DioneStakingCaller) StartBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DioneStaking.contract.Call(opts, out, "startBlock")
	return *ret0, err
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() view returns(uint256)
func (_DioneStaking *DioneStakingSession) StartBlock() (*big.Int, error) {
	return _DioneStaking.Contract.StartBlock(&_DioneStaking.CallOpts)
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() view returns(uint256)
func (_DioneStaking *DioneStakingCallerSession) StartBlock() (*big.Int, error) {
	return _DioneStaking.Contract.StartBlock(&_DioneStaking.CallOpts)
}

// TotalStake is a free data retrieval call binding the contract method 0x8b0e9f3f.
//
// Solidity: function totalStake() view returns(uint256)
func (_DioneStaking *DioneStakingCaller) TotalStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DioneStaking.contract.Call(opts, out, "totalStake")
	return *ret0, err
}

// TotalStake is a free data retrieval call binding the contract method 0x8b0e9f3f.
//
// Solidity: function totalStake() view returns(uint256)
func (_DioneStaking *DioneStakingSession) TotalStake() (*big.Int, error) {
	return _DioneStaking.Contract.TotalStake(&_DioneStaking.CallOpts)
}

// TotalStake is a free data retrieval call binding the contract method 0x8b0e9f3f.
//
// Solidity: function totalStake() view returns(uint256)
func (_DioneStaking *DioneStakingCallerSession) TotalStake() (*big.Int, error) {
	return _DioneStaking.Contract.TotalStake(&_DioneStaking.CallOpts)
}

// IsLegitMiner is a paid mutator transaction binding the contract method 0x3b1175eb.
//
// Solidity: function isLegitMiner(address _minerAddr) returns(bool)
func (_DioneStaking *DioneStakingTransactor) IsLegitMiner(opts *bind.TransactOpts, _minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "isLegitMiner", _minerAddr)
}

// IsLegitMiner is a paid mutator transaction binding the contract method 0x3b1175eb.
//
// Solidity: function isLegitMiner(address _minerAddr) returns(bool)
func (_DioneStaking *DioneStakingSession) IsLegitMiner(_minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.Contract.IsLegitMiner(&_DioneStaking.TransactOpts, _minerAddr)
}

// IsLegitMiner is a paid mutator transaction binding the contract method 0x3b1175eb.
//
// Solidity: function isLegitMiner(address _minerAddr) returns(bool)
func (_DioneStaking *DioneStakingTransactorSession) IsLegitMiner(_minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.Contract.IsLegitMiner(&_DioneStaking.TransactOpts, _minerAddr)
}

// Mine is a paid mutator transaction binding the contract method 0x81923240.
//
// Solidity: function mine(address _minerAddr) returns()
func (_DioneStaking *DioneStakingTransactor) Mine(opts *bind.TransactOpts, _minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "mine", _minerAddr)
}

// Mine is a paid mutator transaction binding the contract method 0x81923240.
//
// Solidity: function mine(address _minerAddr) returns()
func (_DioneStaking *DioneStakingSession) Mine(_minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.Contract.Mine(&_DioneStaking.TransactOpts, _minerAddr)
}

// Mine is a paid mutator transaction binding the contract method 0x81923240.
//
// Solidity: function mine(address _minerAddr) returns()
func (_DioneStaking *DioneStakingTransactorSession) Mine(_minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.Contract.Mine(&_DioneStaking.TransactOpts, _minerAddr)
}

// MineAndStake is a paid mutator transaction binding the contract method 0x407b4547.
//
// Solidity: function mineAndStake(address _minerAddr) returns()
func (_DioneStaking *DioneStakingTransactor) MineAndStake(opts *bind.TransactOpts, _minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "mineAndStake", _minerAddr)
}

// MineAndStake is a paid mutator transaction binding the contract method 0x407b4547.
//
// Solidity: function mineAndStake(address _minerAddr) returns()
func (_DioneStaking *DioneStakingSession) MineAndStake(_minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.Contract.MineAndStake(&_DioneStaking.TransactOpts, _minerAddr)
}

// MineAndStake is a paid mutator transaction binding the contract method 0x407b4547.
//
// Solidity: function mineAndStake(address _minerAddr) returns()
func (_DioneStaking *DioneStakingTransactorSession) MineAndStake(_minerAddr common.Address) (*types.Transaction, error) {
	return _DioneStaking.Contract.MineAndStake(&_DioneStaking.TransactOpts, _minerAddr)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DioneStaking *DioneStakingTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DioneStaking *DioneStakingSession) RenounceOwnership() (*types.Transaction, error) {
	return _DioneStaking.Contract.RenounceOwnership(&_DioneStaking.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DioneStaking *DioneStakingTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DioneStaking.Contract.RenounceOwnership(&_DioneStaking.TransactOpts)
}

// SetMinerReward is a paid mutator transaction binding the contract method 0x816caed5.
//
// Solidity: function setMinerReward(uint256 _minerReward) returns()
func (_DioneStaking *DioneStakingTransactor) SetMinerReward(opts *bind.TransactOpts, _minerReward *big.Int) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "setMinerReward", _minerReward)
}

// SetMinerReward is a paid mutator transaction binding the contract method 0x816caed5.
//
// Solidity: function setMinerReward(uint256 _minerReward) returns()
func (_DioneStaking *DioneStakingSession) SetMinerReward(_minerReward *big.Int) (*types.Transaction, error) {
	return _DioneStaking.Contract.SetMinerReward(&_DioneStaking.TransactOpts, _minerReward)
}

// SetMinerReward is a paid mutator transaction binding the contract method 0x816caed5.
//
// Solidity: function setMinerReward(uint256 _minerReward) returns()
func (_DioneStaking *DioneStakingTransactorSession) SetMinerReward(_minerReward *big.Int) (*types.Transaction, error) {
	return _DioneStaking.Contract.SetMinerReward(&_DioneStaking.TransactOpts, _minerReward)
}

// SetMinimumStake is a paid mutator transaction binding the contract method 0x233e9903.
//
// Solidity: function setMinimumStake(uint256 _minimumStake) returns()
func (_DioneStaking *DioneStakingTransactor) SetMinimumStake(opts *bind.TransactOpts, _minimumStake *big.Int) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "setMinimumStake", _minimumStake)
}

// SetMinimumStake is a paid mutator transaction binding the contract method 0x233e9903.
//
// Solidity: function setMinimumStake(uint256 _minimumStake) returns()
func (_DioneStaking *DioneStakingSession) SetMinimumStake(_minimumStake *big.Int) (*types.Transaction, error) {
	return _DioneStaking.Contract.SetMinimumStake(&_DioneStaking.TransactOpts, _minimumStake)
}

// SetMinimumStake is a paid mutator transaction binding the contract method 0x233e9903.
//
// Solidity: function setMinimumStake(uint256 _minimumStake) returns()
func (_DioneStaking *DioneStakingTransactorSession) SetMinimumStake(_minimumStake *big.Int) (*types.Transaction, error) {
	return _DioneStaking.Contract.SetMinimumStake(&_DioneStaking.TransactOpts, _minimumStake)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) returns()
func (_DioneStaking *DioneStakingTransactor) Stake(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "stake", _amount)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) returns()
func (_DioneStaking *DioneStakingSession) Stake(_amount *big.Int) (*types.Transaction, error) {
	return _DioneStaking.Contract.Stake(&_DioneStaking.TransactOpts, _amount)
}

// Stake is a paid mutator transaction binding the contract method 0xa694fc3a.
//
// Solidity: function stake(uint256 _amount) returns()
func (_DioneStaking *DioneStakingTransactorSession) Stake(_amount *big.Int) (*types.Transaction, error) {
	return _DioneStaking.Contract.Stake(&_DioneStaking.TransactOpts, _amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DioneStaking *DioneStakingTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DioneStaking *DioneStakingSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DioneStaking.Contract.TransferOwnership(&_DioneStaking.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DioneStaking *DioneStakingTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DioneStaking.Contract.TransferOwnership(&_DioneStaking.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_DioneStaking *DioneStakingTransactor) Withdraw(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _DioneStaking.contract.Transact(opts, "withdraw", _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_DioneStaking *DioneStakingSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _DioneStaking.Contract.Withdraw(&_DioneStaking.TransactOpts, _amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 _amount) returns()
func (_DioneStaking *DioneStakingTransactorSession) Withdraw(_amount *big.Int) (*types.Transaction, error) {
	return _DioneStaking.Contract.Withdraw(&_DioneStaking.TransactOpts, _amount)
}

// DioneStakingMineIterator is returned from FilterMine and is used to iterate over the raw logs and unpacked data for Mine events raised by the DioneStaking contract.
type DioneStakingMineIterator struct {
	Event *DioneStakingMine // Event containing the contract specifics and raw log

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
func (it *DioneStakingMineIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneStakingMine)
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
		it.Event = new(DioneStakingMine)
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
func (it *DioneStakingMineIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneStakingMineIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneStakingMine represents a Mine event raised by the DioneStaking contract.
type DioneStakingMine struct {
	Miner       common.Address
	BlockNumber *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterMine is a free log retrieval operation binding the contract event 0xf23a961744a760027f8811c59a0eaef0d29cf965578b17412bcc375b52fa39d1.
//
// Solidity: event Mine(address indexed miner, uint256 blockNumber)
func (_DioneStaking *DioneStakingFilterer) FilterMine(opts *bind.FilterOpts, miner []common.Address) (*DioneStakingMineIterator, error) {

	var minerRule []interface{}
	for _, minerItem := range miner {
		minerRule = append(minerRule, minerItem)
	}

	logs, sub, err := _DioneStaking.contract.FilterLogs(opts, "Mine", minerRule)
	if err != nil {
		return nil, err
	}
	return &DioneStakingMineIterator{contract: _DioneStaking.contract, event: "Mine", logs: logs, sub: sub}, nil
}

// WatchMine is a free log subscription operation binding the contract event 0xf23a961744a760027f8811c59a0eaef0d29cf965578b17412bcc375b52fa39d1.
//
// Solidity: event Mine(address indexed miner, uint256 blockNumber)
func (_DioneStaking *DioneStakingFilterer) WatchMine(opts *bind.WatchOpts, sink chan<- *DioneStakingMine, miner []common.Address) (event.Subscription, error) {

	var minerRule []interface{}
	for _, minerItem := range miner {
		minerRule = append(minerRule, minerItem)
	}

	logs, sub, err := _DioneStaking.contract.WatchLogs(opts, "Mine", minerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneStakingMine)
				if err := _DioneStaking.contract.UnpackLog(event, "Mine", log); err != nil {
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

// ParseMine is a log parse operation binding the contract event 0xf23a961744a760027f8811c59a0eaef0d29cf965578b17412bcc375b52fa39d1.
//
// Solidity: event Mine(address indexed miner, uint256 blockNumber)
func (_DioneStaking *DioneStakingFilterer) ParseMine(log types.Log) (*DioneStakingMine, error) {
	event := new(DioneStakingMine)
	if err := _DioneStaking.contract.UnpackLog(event, "Mine", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DioneStakingOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DioneStaking contract.
type DioneStakingOwnershipTransferredIterator struct {
	Event *DioneStakingOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DioneStakingOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneStakingOwnershipTransferred)
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
		it.Event = new(DioneStakingOwnershipTransferred)
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
func (it *DioneStakingOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneStakingOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneStakingOwnershipTransferred represents a OwnershipTransferred event raised by the DioneStaking contract.
type DioneStakingOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DioneStaking *DioneStakingFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DioneStakingOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DioneStaking.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DioneStakingOwnershipTransferredIterator{contract: _DioneStaking.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DioneStaking *DioneStakingFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DioneStakingOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DioneStaking.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneStakingOwnershipTransferred)
				if err := _DioneStaking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_DioneStaking *DioneStakingFilterer) ParseOwnershipTransferred(log types.Log) (*DioneStakingOwnershipTransferred, error) {
	event := new(DioneStakingOwnershipTransferred)
	if err := _DioneStaking.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DioneStakingStakeIterator is returned from FilterStake and is used to iterate over the raw logs and unpacked data for Stake events raised by the DioneStaking contract.
type DioneStakingStakeIterator struct {
	Event *DioneStakingStake // Event containing the contract specifics and raw log

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
func (it *DioneStakingStakeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneStakingStake)
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
		it.Event = new(DioneStakingStake)
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
func (it *DioneStakingStakeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneStakingStakeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneStakingStake represents a Stake event raised by the DioneStaking contract.
type DioneStakingStake struct {
	Miner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStake is a free log retrieval operation binding the contract event 0xebedb8b3c678666e7f36970bc8f57abf6d8fa2e828c0da91ea5b75bf68ed101a.
//
// Solidity: event Stake(address indexed miner, uint256 amount)
func (_DioneStaking *DioneStakingFilterer) FilterStake(opts *bind.FilterOpts, miner []common.Address) (*DioneStakingStakeIterator, error) {

	var minerRule []interface{}
	for _, minerItem := range miner {
		minerRule = append(minerRule, minerItem)
	}

	logs, sub, err := _DioneStaking.contract.FilterLogs(opts, "Stake", minerRule)
	if err != nil {
		return nil, err
	}
	return &DioneStakingStakeIterator{contract: _DioneStaking.contract, event: "Stake", logs: logs, sub: sub}, nil
}

// WatchStake is a free log subscription operation binding the contract event 0xebedb8b3c678666e7f36970bc8f57abf6d8fa2e828c0da91ea5b75bf68ed101a.
//
// Solidity: event Stake(address indexed miner, uint256 amount)
func (_DioneStaking *DioneStakingFilterer) WatchStake(opts *bind.WatchOpts, sink chan<- *DioneStakingStake, miner []common.Address) (event.Subscription, error) {

	var minerRule []interface{}
	for _, minerItem := range miner {
		minerRule = append(minerRule, minerItem)
	}

	logs, sub, err := _DioneStaking.contract.WatchLogs(opts, "Stake", minerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneStakingStake)
				if err := _DioneStaking.contract.UnpackLog(event, "Stake", log); err != nil {
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

// ParseStake is a log parse operation binding the contract event 0xebedb8b3c678666e7f36970bc8f57abf6d8fa2e828c0da91ea5b75bf68ed101a.
//
// Solidity: event Stake(address indexed miner, uint256 amount)
func (_DioneStaking *DioneStakingFilterer) ParseStake(log types.Log) (*DioneStakingStake, error) {
	event := new(DioneStakingStake)
	if err := _DioneStaking.contract.UnpackLog(event, "Stake", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DioneStakingWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the DioneStaking contract.
type DioneStakingWithdrawIterator struct {
	Event *DioneStakingWithdraw // Event containing the contract specifics and raw log

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
func (it *DioneStakingWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DioneStakingWithdraw)
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
		it.Event = new(DioneStakingWithdraw)
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
func (it *DioneStakingWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DioneStakingWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DioneStakingWithdraw represents a Withdraw event raised by the DioneStaking contract.
type DioneStakingWithdraw struct {
	Miner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed miner, uint256 amount)
func (_DioneStaking *DioneStakingFilterer) FilterWithdraw(opts *bind.FilterOpts, miner []common.Address) (*DioneStakingWithdrawIterator, error) {

	var minerRule []interface{}
	for _, minerItem := range miner {
		minerRule = append(minerRule, minerItem)
	}

	logs, sub, err := _DioneStaking.contract.FilterLogs(opts, "Withdraw", minerRule)
	if err != nil {
		return nil, err
	}
	return &DioneStakingWithdrawIterator{contract: _DioneStaking.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed miner, uint256 amount)
func (_DioneStaking *DioneStakingFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *DioneStakingWithdraw, miner []common.Address) (event.Subscription, error) {

	var minerRule []interface{}
	for _, minerItem := range miner {
		minerRule = append(minerRule, minerItem)
	}

	logs, sub, err := _DioneStaking.contract.WatchLogs(opts, "Withdraw", minerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DioneStakingWithdraw)
				if err := _DioneStaking.contract.UnpackLog(event, "Withdraw", log); err != nil {
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

// ParseWithdraw is a log parse operation binding the contract event 0x884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a9424364.
//
// Solidity: event Withdraw(address indexed miner, uint256 amount)
func (_DioneStaking *DioneStakingFilterer) ParseWithdraw(log types.Log) (*DioneStakingWithdraw, error) {
	event := new(DioneStakingWithdraw)
	if err := _DioneStaking.contract.UnpackLog(event, "Withdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}
