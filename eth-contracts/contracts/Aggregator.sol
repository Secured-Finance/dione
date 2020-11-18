//SPDX-License-Identifier: MIT
pragma solidity >= 0.5.0 < 0.7.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

interface IDioneStaking {
    function mine(address _minerAddr) external;
    function mineAndStake(address _minerAddr) external;
    function isLegitMiner(address _minerAddr) external returns (bool);
}

interface IMediator {
    function _receiveDataCallback(uint256 reqID, string memory data) external;
}

contract Aggregator is Ownable, ReentrancyGuard {
  IDioneStaking public dioneStaking;

  // Set DioneStaking contract. Can only be called by the owner.
  function setDioneStaking(IDioneStaking _dioneStaking) public onlyOwner {
      dioneStaking = _dioneStaking;
  }

  function collectData(uint256 reqID, string memory data, IMediator callbackAddress) public nonReentrant {
    callbackAddress._receiveDataCallback(reqID, data);
  }
}
