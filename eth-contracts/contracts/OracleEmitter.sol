// SPDX-License-Identifier: MIT
pragma solidity >=0.4.21 <0.7.0;

contract OracleEmitter {
  uint256 requestCounter;

  event NewOracleRequest(
    uint8 originChain,
    string requestType,
    string requestParams,
    address callbackAddress,
    bytes4 callbackMethodID,
    uint256 requestID
  );

  function requestOracles(uint8 originChain, string memory requestType, string memory requestParams, address callbackAddress, bytes4 callbackMethodID) public returns (uint256) {
    requestCounter++;
    emit NewOracleRequest(originChain, requestType, requestParams, callbackAddress, callbackMethodID, requestCounter);
    return requestCounter;
  }
}
