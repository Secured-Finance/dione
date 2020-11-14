//SPDX-License-Identifier: MIT
pragma solidity >= 0.5.0 < 0.7.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";

interface IDioneStaking {
    function mine(address _minerAddr) external;
    function mineAndStake(address _minerAddr) external;
    function isLegitMiner(address _minerAddr) external returns (bool);
}

contract Aggregator is Ownable, ReentrancyGuard {
  IDioneStaking public dioneStaking;

  // Set DioneStaking contract. Can only be called by the owner.
  function setDioneStaking(IDioneStaking _dioneStaking) public onlyOwner {
      dioneStaking = _dioneStaking;
  }

  function collectData(uint256 reqID, string memory data, address callbackAddress, bytes4 callbackMethodID) public nonReentrant {
    require(dioneStaking.isLegitMiner(msg.sender));
    (bool success,) = callbackAddress.call(abi.encode(callbackMethodID, reqID, data));
    require(success);
  }
}
