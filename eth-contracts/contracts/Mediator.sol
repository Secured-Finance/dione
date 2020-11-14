//SPDX-License-Identifier: MIT
pragma solidity >= 0.5.0 < 0.7.0;

import "./vendor/Ownable.sol";
import "./OracleEmitter.sol";

contract Mediator is Ownable {
  event DataReceived(
    uint256 reqID,
    string data
  );

  OracleEmitter oracleEmitter;
  address aggregator;

  constructor(address _oracleEmitter, address _aggregator) public Ownable() {
    oracleEmitter = OracleEmitter(_oracleEmitter);
    aggregator = _aggregator;
  }

  function request(uint8 memory originChain, uint8 memory requestType, string memory requestParams) public returns (uint256) {
    return oracleEmitter.requestOracles(originChain, requestType, requestParams, address(this), bytes4(keccak256("_receiveDataCallback(uint256, string)")));
  }

  function _receiveDataCallback(uint256 reqID, string memory data) public {
    require(msg.sender == aggregator);
    emit DataReceived(reqID, data);
  }
}
