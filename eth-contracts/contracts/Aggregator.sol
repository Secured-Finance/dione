//SPDX-License-Identifier: MIT
pragma solidity >= 0.5.0 < 0.7.0;

import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "./DioneStaking.sol";

interface IMediator {
    function _receiveDataCallback(uint256 reqID, string memory data) external;
}

contract Aggregator is Ownable, ReentrancyGuard {
  DioneStaking public dioneStaking;

  constructor(DioneStaking _dioneStaking) public {
      dioneStaking = _dioneStaking;
  }

  function collectData(uint256 reqID, string memory data, IMediator callbackAddress) public {
    dioneStaking.mine(msg.sender);
    callbackAddress._receiveDataCallback(reqID, data);
  }
}
