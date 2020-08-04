//SPDX-License-Identifier: MIT
pragma solidity >= 0.5.0 < 0.7.0;

import "./Whitelist.sol";

contract Aggregator {
  Whitelist whitelist;

  constructor(Whitelist _whitelist) public {
    whitelist = _whitelist;
  }

  function collectData(uint256 reqID, string memory data, address callbackAddress, bytes4 callbackMethodID) public {
    address[] memory nodes = whitelist.getNodes();
    bool nodeInTheList = false;
    for (uint256 i = 0; i < nodes.length-1; i++) {
      if (nodes[i] == msg.sender) {
        nodeInTheList = true;
      }
    }
    require(nodeInTheList);
    (bool success,) = callbackAddress.call(abi.encode(callbackMethodID, reqID, data));
    require(success);
  }
}
