//SPDX-License-Identifier: MIT
pragma solidity >= 0.5.0 < 0.7.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./DioneOracle.sol";

contract Mediator is Ownable {
  event DataReceived(
    uint256 reqID,
    string data
  );

  DioneOracle dioneOracle;
  address aggregator;

  constructor(DioneOracle _dioneOracle) public Ownable() {
    dioneOracle = _dioneOracle;
  }

  function request(uint8 originChain, string memory requestType, string memory requestParams) public returns (uint256) {
    return dioneOracle.requestOracles(originChain, requestType, requestParams, address(this), bytes4(keccak256("_receiveDataCallback(uint256, string)")));
  }

  function _receiveDataCallback(uint256 reqID, string memory data) public {
    require(msg.sender == aggregator);
    emit DataReceived(reqID, data);
  }
}
