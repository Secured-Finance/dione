pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "./interfaces/IDioneStaking.sol";

contract DioneDispute {
    IDioneStaking public dioneStaking;
    uint256 public voteWindowTime;

    // Minimum amount of DIONE tokens required to vote to dispute
    uint256 public minStake;

    struct Dispute {
        bytes32 dhash; // id of dispute - keccak256(_miner,_requestId,_timestamp)
        int256 sum; // vote measure (for/against this dispute)
        bool finished; // dispute was finished (closed) or not
        bool disputeResult; // true - dispute had basis, false - dispute was false
        address miner; // the miner against whom the dispute
        address disputeInitiator; // the miner who started the dispute
        uint256 timestamp; // dispute creation timestamp
        address[] voted; // map of miners who vote for/against
    }

    modifier onlyExistingDispute(bytes32 _dhash) {
        require(disputes[_dhash].disputeInitiator != address(0), "dispute doesn't exist");
        _;
    } 

    mapping(bytes32 => Dispute) disputes;

    event NewDispute(bytes32 dhash, uint256 requestID, address indexed miner, address indexed disputeInitiator);
    event NewVote(bytes32 dhash, address indexed votedMiner);
    event DisputeFinished(bytes32 dhash, bool status);

    constructor(IDioneStaking _dioneStaking, uint256 _voteWindowTime, uint256 _minStake) {
        dioneStaking = _dioneStaking;
        voteWindowTime = _voteWindowTime;
        minStake = _minStake;
    }

    function beginDispute(address miner, uint256 requestID) public {
        bytes32 dhash = keccak256(abi.encodePacked(miner, requestID));
        require(disputes[dhash].miner == address(0), "dispute already exists");
        require(dioneStaking.isMiner(msg.sender), "caller isn't dione miner");
        Dispute storage dispute = disputes[dhash];
        dispute.dhash = dhash;
        dispute.sum = 0;
        dispute.finished = false;
        dispute.disputeResult = false;
        dispute.miner = miner;
        dispute.timestamp = block.timestamp;
        dispute.disputeInitiator = msg.sender;

        disputes[dhash] = dispute;

        emit NewDispute(dhash, requestID, miner, msg.sender);
    }

    function vote(bytes32 dhash, bool voteStatus) public onlyExistingDispute(dhash) {
        Dispute memory dispute = disputes[dhash];
        require(dispute.finished == false, "dispute already finished");
        require(msg.sender != dispute.disputeInitiator, "dispute initiator isn't allowed to vote");
        require(msg.sender != dispute.miner, "the miner against whom dispute has beginned isn't allowed to vote");
        require(dioneStaking.isMiner(msg.sender), "caller isn't dione miner");
        require(dioneStaking.minerStake(msg.sender) >= minStake, "miner doesn't have minimum stake to vote");
        uint256 stake = dioneStaking.minerStake(msg.sender);
        if (voteStatus) {
            disputes[dhash].sum = disputes[dhash].sum + int256(stake);
        } else {
            disputes[dhash].sum = disputes[dhash].sum - int256(stake);
        }
        disputes[dhash].voted.push(msg.sender);

        emit NewVote(dhash, msg.sender);
    }

    function finishDispute(bytes32 dhash) public onlyExistingDispute(dhash) {
        Dispute memory dispute = disputes[dhash];
        require((block.timestamp - dispute.timestamp) >= voteWindowTime, "vote window hasn't passed yet");
        require(dispute.finished == false, "dispute already finished");
        require(dispute.disputeInitiator == msg.sender, "only dispute initiator can call this function");
        if (dispute.sum <= 0) {
            disputes[dhash].disputeResult = false;
        } else {
            disputes[dhash].disputeResult = true;
            disputes[dhash].voted.push(msg.sender);
            dioneStaking.slashMiner(dispute.miner, disputes[dhash].voted);
        }

        disputes[dhash].finished = true;

        emit DisputeFinished(dhash, disputes[dhash].disputeResult);
    }
}
