// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/utils/math/SafeMath.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

import "./interfaces/IDioneStaking.sol";

contract DioneOracle {
  using SafeMath for uint256;
  
  // Global counter of oracle requests, works as an index in mapping structures
  uint256 private requestCounter = 0;
  // Maximum time for computing oracle request
  uint256 constant public MAXIMUM_DELAY = 5 minutes;
  // Dione staking contract
  IDioneStaking public dioneStaking;
  // Minimum amount of DIONE tokens required to vote against miner result
  uint256 public minimumDisputeFee = 100**18;

  struct OracleRequest {
    address requestSender;
    uint8 originChain; // origin blockchain for request
    string requestType; // rpc call type
    string requestParams; // rpc call params
    address callbackAddress; // callback address for users request contract
    bytes4 callbackMethodID; // method for users request contract
    uint256 reqID; // request counter 
    uint256 deadline;
    bytes data;
  }

  mapping(uint256 => OracleRequest) private pendingRequests;

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
    uint256 reqID,
    bytes data
  );

  modifier onlyPendingRequest(uint256 _reqID) {
    require(pendingRequests[_reqID].requestSender != address(0), "This request is not pending");
    _;
  }

  constructor(IDioneStaking _dioneStaking) {
      dioneStaking = _dioneStaking;
  }

  function requestOracles(uint8 _originChain, string memory _requestType, string memory _requestParams, address _callbackAddress, bytes4 _callbackMethodID) public returns (uint256) {
    requestCounter += 1;
    require(pendingRequests[requestCounter].reqID == 0, "This counter is not unique");
    uint256 requestDeadline = block.timestamp.add(MAXIMUM_DELAY);
    pendingRequests[requestCounter] = OracleRequest({
      requestSender: msg.sender,
      originChain: _originChain,
      requestType: _requestType,
      requestParams: _requestParams,
      callbackAddress: _callbackAddress,
      callbackMethodID: _callbackMethodID,
      reqID: requestCounter,
      deadline: requestDeadline,
      data: new bytes(0)
    });

    emit NewOracleRequest(_originChain, _requestType, _requestParams, _callbackAddress, _callbackMethodID, requestCounter, requestDeadline);
    return requestCounter;
  }

  function cancelOracleRequest(uint256 _reqID) public onlyPendingRequest(_reqID) {
    require(msg.sender == pendingRequests[_reqID].requestSender, "you aren't request sender");

    delete pendingRequests[_reqID];
    emit CancelOracleRequest(_reqID);
  }

  function submitOracleRequest(uint256 _reqID, bytes memory _data) public onlyPendingRequest(_reqID) returns (bool) {
    delete pendingRequests[_reqID];
    dioneStaking.mine(msg.sender);
    (bool success, ) = pendingRequests[_reqID].callbackAddress.call(abi.encodeWithSelector(pendingRequests[_reqID].callbackMethodID, _reqID, _data));
    require(success == true, "cannot call callback method");
    emit SubmittedOracleRequest(_reqID, _data);
    return true;
  }
}
