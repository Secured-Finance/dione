// SPDX-License-Identifier: MIT
pragma solidity >=0.4.21 <0.7.0;

contract OracleEmitter {
  uint256 requestCounter;

  event NewOracleRequest(
    string requestType,
    address callbackAddress,
    bytes4 callbackMethodID,
    uint256 requestID
  );

  function requestOracles(string memory requestType, address callbackAddress, bytes4 callbackMethodID) public returns (uint256) {
    requestCounter++;
    emit NewOracleRequest(requestType, callbackAddress, callbackMethodID, requestCounter);
    return requestCounter;
  }
}
