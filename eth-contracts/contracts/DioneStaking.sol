pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
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
    // DioneOracle contract address.
    address public dioneOracleAddress;
    // Dispute contract address.
    address public disputeContractAddr;
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
    event RewardChanged(uint256 oldValue, uint256 newValue);
    event MinimumStakeChanged(uint256 oldValue, uint256 newValue);

    constructor(
        DioneToken _dione,
        uint256 _minerReward,
        uint256 _startBlock,
        uint256 _minimumStake
    ) {
        dione = _dione;
        minerReward = _minerReward;
        startBlock = _startBlock;
        minimumStake = _minimumStake;
    }

    // Mine new dione oracle task, only can be executed by oracle contract
    function mine(address _minerAddr) public nonReentrant {
        require(msg.sender == dioneOracleAddress, "not oracle contract");
        MinerInfo storage miner = minerInfo[_minerAddr];
        dione.mint(_minerAddr, minerReward);
        miner.lastRewardBlock = block.number;
        emit Mine(_minerAddr, block.number);
    }

    // Mine new dione oracle task and stake miner reward, only can be executed by oracle contract
    function mineAndStake(address _minerAddr) public nonReentrant {
        require(msg.sender == dioneOracleAddress, "not oracle contract");
        MinerInfo storage miner = minerInfo[_minerAddr];
        dione.mint(address(this), minerReward);
        _totalStake = _totalStake.add(minerReward);
        miner.amount = miner.amount.add(minerReward);
        miner.lastRewardBlock = block.number;
        emit Mine(_minerAddr, block.number);
    }

    // Deposit DIONE tokens to mine on dione network
    function stake(uint256 _amount) public nonReentrant {
        require(_amount > 0, "cannot stake zero");
        require(_amount >= minimumStake, "actual stake amount is less than minimum stake amount");
        MinerInfo storage miner = minerInfo[msg.sender];
        dione.transferFrom(address(msg.sender), address(this), _amount);
        _totalStake = _totalStake.add(_amount);
        miner.amount = miner.amount.add(_amount);
        if (miner.firstStakeBlock == 0) {
            miner.firstStakeBlock = block.number > startBlock ? block.number : startBlock;
        }
        emit Stake(msg.sender, _amount);
    }

    // Withdraw DIONE tokens from DioneStaking
    function withdraw(uint256 _amount) public nonReentrant {
        MinerInfo storage miner = minerInfo[msg.sender];
        require(miner.amount >= _amount, "withdraw: not enough tokens");
        require(_amount > 0, "cannot withdraw zero");
        _totalStake = _totalStake.sub(_amount);
        miner.amount = miner.amount.sub(_amount);
        dione.transfer(address(msg.sender), _amount);
        emit Withdraw(msg.sender, _amount);
    }

    // Returns total amount of DIONE tokens in PoS mining
    function totalStake() external view returns (uint256) {
        return _totalStake;
    }

    function minerStake(address _minerAddr) external view returns (uint256) {
       return minerInfo[_minerAddr].amount;
    }

    // Update miner reward in DIONE tokens, only can be executed by owner of the contract
    function setMinerReward(uint256 _minerReward) public onlyOwner {
        require(_minerReward > 0, "reward must not be zero");
        uint256 oldRewardValue = minerReward;
        minerReward = _minerReward;
        emit RewardChanged(oldRewardValue, _minerReward);
    }

    function isMiner(address _minerAddr) public view returns (bool) {
        return minerInfo[_minerAddr].amount >= minimumStake;
    }

    // Update minimum stake in DIONE tokens for miners, only can be executed by owner of the contract
    function setMinimumStake(uint256 _minimumStake) public onlyOwner {
        require(_minimumStake > 0, "minimum stake must not be zero");
        uint256 oldValue = minimumStake;
        minimumStake = _minimumStake;
        emit MinimumStakeChanged(oldValue, _minimumStake);
    }

    function setOracleContractAddress(address _addr) public onlyOwner {
        require(_addr != address(0), "address must not be zero");
        dioneOracleAddress = _addr;
    }

    function setDisputeContractAddress(address _addr) public onlyOwner {
        require(_addr != address(0), "address must not be zero");
        disputeContractAddr = _addr;
    }

    function slashMiner(address miner, address[] memory receipentMiners) public {
        require(msg.sender == disputeContractAddr, "caller is not the dispute contract");
        require(miner != address(0), "slashing address must not be zero");
        require(isMiner(miner), "slashing address isn't dione miner");
        
        uint256 share = minerInfo[miner].amount.div(receipentMiners.length);


        for (uint8 i = 0; i < receipentMiners.length; i++) {
            require(receipentMiners[i] != miner, "receipent address must not be slashing address");
            require(isMiner(receipentMiners[i]), "receipent address isn't dione miner");
            minerInfo[miner].amount = minerInfo[miner].amount.sub(share);
            minerInfo[receipentMiners[i]].amount = minerInfo[receipentMiners[i]].amount.add(share);
        }
    }
}
