pragma solidity 0.6.12;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/math/SafeMath.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./DioneToken.sol";

// DioneStaking is the main contract for Dione oracle network Proof-of-Stake mechanism
//
// Note that it's ownable and the owner has power over changing miner rewards and minimum DIONE stake
// The ownership of the contract would be transfered to 24 hours Timelock
contract DioneStaking is Ownable, ReentrancyGuard {
    using SafeMath for uint256;

    // MinerInfo contains total DIONEs staked by miner, how much tasks been already computed
    // and timestamp of first deposit
    struct MinerInfo {
        uint256 amount; // How many DIONE tokens was staked by the miner.
        uint256 firstStakeBlock; // First block for miner DIONE tokens reward.
        uint256 lastRewardBlock; // Last block for miner DIONE tokens reward
    }

    DioneToken public dione;
    // Agregator contract address.
    address public aggregatorAddr;
    // Miner rewards in DIONE tokens.
    uint256 public minerReward;
    // The block number when DIONE mining starts.
    uint256 public startBlock;
    // Minimum amount of staked DIONE tokens required to start mining
    uint256 public minimumStake;
    // Total amount of DIONE tokens staked
    uint256 private _totalStake;

    // Info of each miner that stakes DIONE tokens.
    mapping (address => MinerInfo) public minerInfo;

    event Stake(address indexed miner, uint256 amount);
    event Withdraw(address indexed miner, uint256 amount);
    event Mine(address indexed miner, uint256 blockNumber);

    constructor(
        DioneToken _dione,
        address _aggregatorAddr,
        uint256 _minerReward,
        uint256 _startBlock,
        uint256 _minimumStake
    ) public {
        dione = _dione;
        aggregatorAddr = _aggregatorAddr;
        minerReward = _minerReward;
        startBlock = _startBlock;
        minimumStake = _minimumStake;
    }

    // Mine new dione oracle task, only can be executed by aggregator contract
    function mine(address _minerAddr) public nonReentrant {
        require(msg.sender == aggregatorAddr, "not aggregator contract");
        MinerInfo storage miner = minerInfo[_minerAddr];
        require(miner.amount >= minimumStake);
        dione.mint(_minerAddr, minerReward);
        miner.lastRewardBlock = block.number;
        emit Mine(_minerAddr, block.number);
    }

    // Mine new dione oracle task and stake miner reward, only can be executed by aggregator contract
    function mineAndStake(address _minerAddr) public nonReentrant {
        require(msg.sender == aggregatorAddr, "not aggregator contract");
        MinerInfo storage miner = minerInfo[_minerAddr];
        require(miner.amount >= minimumStake);
        dione.mint(address(this), minerReward);
        _totalStake = _totalStake.add(minerReward);
        miner.amount = miner.amount.add(minerReward);
        miner.lastRewardBlock = block.number;
        emit Mine(_minerAddr, block.number);
    }

    // Deposit DIONE tokens to mine on dione network
    function stake(uint256 _amount) public nonReentrant {
        require(_amount > 0, "Cannot stake 0");
        MinerInfo storage miner = minerInfo[msg.sender];
        _totalStake = _totalStake.add(_amount);
        miner.amount = miner.amount.add(_amount);
        dione.transferFrom(address(msg.sender), address(this), _amount);
        if (miner.firstStakeBlock == 0 && miner.amount >= minimumStake) {
            miner.firstStakeBlock = block.number > startBlock ? block.number : startBlock;
        }
        emit Stake(msg.sender, _amount);
    }

    // Withdraw DIONE tokens from DioneStaking
    function withdraw(uint256 _amount) public nonReentrant {
        MinerInfo storage miner = minerInfo[msg.sender];
        require(miner.amount >= _amount, "withdraw: not enough tokens");
        if(_amount > 0) {
            _totalStake = _totalStake.sub(_amount);
            miner.amount = miner.amount.sub(_amount);
            dione.transfer(address(msg.sender), _amount);
        }
        emit Withdraw(msg.sender, _amount);
    }

    // Returns total amount of DIONE tokens in PoS mining
    function totalStake() external view returns (uint256) {
        return _totalStake;
    }

    function minerStake(address _minerAddr) external view returns (uint256) {
        MinerInfo storage miner = minerInfo[_minerAddr];
        return miner.amount;
    }

    // Update miner reward in DIONE tokens, only can be executed by owner of the contract
    function setMinerReward(uint256 _minerReward) public onlyOwner {
        require(_minerReward > 0, "!minerReward-0");
        minerReward = _minerReward;
    }

    // Update minimum stake in DIONE tokens for miners, only can be executed by owner of the contract
    function setMinimumStake(uint256 _minimumStake) public onlyOwner {
        require(_minimumStake > 0, "!minerReward-0");
        minimumStake = _minimumStake;
    }
}