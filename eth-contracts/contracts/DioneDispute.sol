pragma solidity ^0.6.12;

import "@openzeppelin/contracts/math/SafeMath.sol";
import "./interfaces/IDioneStaking.sol";

contract DioneDispute {
    using SafeMath for uint256;

    IDioneStaking public dioneStaking;

    struct Dispute {
        bytes32 dhash; // id of dispute - keccak256(_miner,_requestId,_timestamp)
        uint256 sum; // vote measure (for/against this dispute)
        bool finished; // dispute was finished (closed) or not
        bool disputeResult; // true - dispute had basis, false - dispute was false
        address miner; // the miner against whom the dispute
        address disputeInitiator; // the miner who started the dispute
        uint256 timestamp; // dispute creation timestamp
        address[] voted; // map of miners who vote for/against
    }

    mapping(bytes32 => Dispute) disputes;

    event NewDispute(bytes32 dhash, uint256 requestID, address indexed miner, address indexed disputeInitiator);
    event NewVote(bytes32 dhash, address indexed votedMiner);
    event DisputeFinished(bytes32 dhash, bool status);

    constructor(IDioneStaking _dioneStaking) public {
        dioneStaking = _dioneStaking;
    }

    function beginDispute(address miner, uint256 requestID) public {
        bytes32 dhash = keccak256(abi.encodePacked(miner, requestID, now));
        require(disputes[dhash].dhash.length != 0, "dispute already exists");
        Dispute storage dispute = disputes[dhash];
        dispute.dhash = dhash;
        dispute.sum = 0;
        dispute.finished = false;
        dispute.disputeResult = false;
        dispute.miner = miner;
        dispute.timestamp = now;
        dispute.disputeInitiator = msg.sender;

        disputes[dhash] = dispute;

        emit NewDispute(dhash, requestID, miner, msg.sender);
    }

    function vote(bytes32 dhash, bool voteStatus) public {
        require(disputes[dhash].dhash.length == 0, "dispute doesn't exist");
        Dispute storage dispute = disputes[dhash];
        require(dispute.finished == false, "dispute already finished");
        require(dioneStaking.isMiner(msg.sender), "caller isn't dione miner");
        uint256 stake = dioneStaking.minerStake(msg.sender);
        if (voteStatus) {
            dispute.sum.sub(stake);
        } else {
            dispute.sum.add(stake);
        }
        dispute.voted.push(msg.sender);

        emit NewVote(dhash, msg.sender);
    }

    function finishDispute(bytes32 dhash) public {
        require(disputes[dhash].dhash.length == 0, "dispute doesn't exist");
        Dispute storage dispute = disputes[dhash];
        require((now - dispute.timestamp) >= 2 hours, "vote window must be two hours");
        require(dispute.finished == false, "dispute already finished");
        require(dispute.disputeInitiator == msg.sender, "only dispute initiator can call this function");
        if (dispute.sum < 0) {
            dispute.disputeResult = false;
        } else {
            dispute.disputeResult = true;
            dispute.voted.push(msg.sender);
            dioneStaking.slashMiner(dispute.miner, dispute.voted);
        }

        dispute.finished = true;

        emit DisputeFinished(dhash, dispute.disputeResult);
    }
}
