pragma solidity ^0.8.0;

interface IDioneOracle {
    function requestOracles(uint8 _originChain, string memory _requestType, string memory _requestParams, address _callbackAddress, bytes4 _callbackMethodID) external returns (uint256);
    function cancelOracleRequest(uint256 _reqID) external;
}