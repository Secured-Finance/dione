pragma solidity ^0.8.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "./interfaces/IDioneOracle.sol";

contract Mediator {
  event DataReceived(
    uint256 reqID,
    string data
  );

  IDioneOracle dioneOracle;

  constructor(IDioneOracle _dioneOracle) {
    dioneOracle = _dioneOracle;
  }

  function request(uint8 originChain, string memory requestType, string memory requestParams) public returns (uint256) {
    return dioneOracle.requestOracles(originChain, requestType, requestParams, address(this), bytes4(keccak256("_receiveDataCallback(uint256, string)")));
  }

  function _receiveDataCallback(uint256 reqID, string memory data) public {
    require(msg.sender == address(dioneOracle));
    emit DataReceived(reqID, data);
  }
}