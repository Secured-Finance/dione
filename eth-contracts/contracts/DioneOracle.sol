// SPDX-License-Identifier: MIT
pragma solidity >=0.6.12;

import "@openzeppelin/contracts/math/SafeMath.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

import "./interfaces/IDioneStaking.sol";

contract DioneOracle is Ownable {
  using SafeMath for uint256;
  
  // Global counter of oracle requests, works as an index in mapping structures
  uint256 private requestCounter;
  // Maximum time for computing oracle request
  uint256 constant public MAXIMUM_DELAY = 5 minutes;
  uint256 constant public FINALIZATION = 3 minutes;
  // Dione staking contract
  IDioneStaking public dioneStaking;
  // Minimum amount of DIONE tokens required to vote against miner result
  uint256 public minimumDisputeFee = 100**18;

  struct OracleRequest {
    uint8 originChain; // origin blockchain for request
    string requestType; // rpc call type
    string requestParams; // rpc call params
    address callbackAddress; // callback address for users request contract
    bytes4 callbackMethodID; // method for users request contract
    uint256 reqID; // request counter 
    uint256 deadline;
    bytes data;
  }

  struct MinedTask {
    uint32 disputes;
    bytes32 taskHash;
    address miner;
    uint256 deadline;
  }

  mapping(uint256 => bytes32) private pendingRequests;
  mapping(uint256 => MinedTask) private nonFinTasks;
  mapping(address => bool) private activeNodes;

  event NewOracleRequest(
    uint8 originChain,
    string requestType,
    string requestParams,
    address callbackAddress,
    bytes4 callbackMethodID,
    uint256 reqID,
    uint256 deadline
  );

  event CancelOracleRequest(
    uint256 reqID
  );

  event SubmittedOracleRequest(
    string requestParams,
    address callbackAddress,
    bytes4 callbackMethodID,
    uint256 reqID,
    uint256 deadline,
    bytes data
  );

  modifier onlyPendingRequest(uint256 _reqID) {
    require(pendingRequests[_reqID] != 0, "This request is not pending");
    _;
  }

  modifier onlyActiveNode() {
    require(activeNodes[msg.sender], "Not an active miner");
    _;
  }

  constructor(IDioneStaking _dioneStaking) public {
      dioneStaking = _dioneStaking;
  }

  function requestOracles(uint8 _originChain, string memory _requestType, string memory _requestParams, address _callbackAddress, bytes4 _callbackMethodID) public returns (uint256) {
    requestCounter += 1;
    require(pendingRequests[requestCounter] == 0, "This counter is not unique");
    uint256 requestDeadline = now.add(MAXIMUM_DELAY);
    pendingRequests[requestCounter] = keccak256(abi.encodePacked(_requestParams, _callbackAddress, _callbackMethodID, requestCounter, requestDeadline));

    emit NewOracleRequest(_originChain, _requestType, _requestParams, _callbackAddress, _callbackMethodID, requestCounter, requestDeadline);
    return requestCounter;
  }

  function cancelOracleRequest(string memory _requestParams, bytes4 _callbackMethodID, uint256 _reqID, uint256 _requestDeadline) public {
    bytes32 requestHash = keccak256(abi.encodePacked(_requestParams, msg.sender, _callbackMethodID, _reqID, _requestDeadline));
    require(requestHash == pendingRequests[_reqID], "Request hash do not match it's origin");
    require(_requestDeadline <= now, "Request didn't reached it's deadline");

    delete pendingRequests[_reqID];
    emit CancelOracleRequest(_reqID);
  }

  function submitOracleRequest(string memory _requestParams, address _callbackAddress, bytes4 _callbackMethodID, uint256 _reqID, uint256 _requestDeadline, bytes memory _data) public onlyPendingRequest(_reqID) onlyActiveNode returns (bool) {
    bytes32 requestHash = keccak256(abi.encodePacked(_requestParams, _callbackAddress, _callbackMethodID, _reqID, _requestDeadline));
    require(pendingRequests[_reqID] == requestHash, "Params do not match request ID");
    delete pendingRequests[_reqID];
    dioneStaking.mine(msg.sender);
    (bool success, ) = _callbackAddress.call(abi.encodeWithSelector(_callbackMethodID, _reqID, _data));
    emit SubmittedOracleRequest(_requestParams, _callbackAddress, _callbackMethodID, _reqID, _requestDeadline, _data);
    MinedTask storage task = nonFinTasks[_reqID];
    task.miner = msg.sender;
    task.taskHash = requestHash;
    return success;
  }
}
