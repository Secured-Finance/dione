pragma solidity ^0.6.12;

import "@openzeppelin/contracts/math/SafeMath.sol";
import "./interfaces/IDioneStaking.sol";

contract DioneDispute {
    using SafeMath for uint256;

    IDioneStaking public dioneStaking;

    struct Dispute {
        bytes32 dhash; // id of dispute - keccak256(_miner,_requestId,_timestamp)
        int256 sum; // vote measure (for/against this dispute)
        bool finished; // dispute was finished (closed) or not
        bool disputeResult; // true - dispute had basis, false - dispute was false
        address miner; // the miner against whom the dispute
        address disputeInitiator; // the miner who started the dispute
        mapping(address => bool) voted; // map of miners who vote for/against
    }

    mapping(bytes32 => Dispute) disputes;

    event NewDispute(bytes32 dhash, address indexed miner, address indexed disputeInitiator);
    event NewVote(bytes32 dhash, address indexed votedMiner);
    event DisputeFinished(bytes32 dhash, bool status);

    constructor(IDioneStaking _dioneStaking) public {
        dioneStaking = _dioneStaking;
    }

    function beginDispute(address miner, uint256 requestID) public {
        bytes32 dhash = keccak256(miner, requestID, now);
        Dispute dispute = Dispute(
            {
                dhash: dhash,
                sum: 0,
                finished: false,
                disputeResult: true,
                miner: miner,
                disputeInitiator: msg.sender,
            }
        );

        disputes[dhash] = dispute;

        emit NewDispute(dhash, miner, msg.sender);
    }

    function vote(bytes32 dhash, bool voteStatus) public {
        require(disputes[dhash]);
        Dispute storage dispute = disputes[dhash];
        require(dispute.finished == false);
        require(dioneStaking.isMiner(msg.sender));
        int256 stake = dioneStaking.minerStake(msg.sender);
        if (voteStatus) {
            dispute.sum.sub(stake);
        } else {
            dispute.sum.add(stake);
        }
        dispute.voted[msg.sender] = voteStatus;

        emit NewVote(dhash, msg.sender);
    }

    function finishDispute(bytes32 dhash) public {
        require(disputes[dhash]);
        Dispute storage dispute = disputes[dhash];
        require(dispute.finished == false);
        require(dispute.disputeInitiator == msg.sender);
        if (dispute.sum < 0) {
            dispute.disputeResult = false;
        } else {
            dispute.disputeResult = true;
        }

        emit DisputeFinished(dhash, dispute.disputeResult);
    }
}
