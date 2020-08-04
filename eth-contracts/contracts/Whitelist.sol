//SPDX-License-Identifier: MIT
pragma solidity >= 0.5.0 < 0.7.0;

import "./vendor/Ownable.sol";

contract Whitelist is Ownable {
  address[] public nodes;

  constructor() public Ownable() {
  }

  function addNode(address nodeAccount) public onlyOwner {
    nodes.push(nodeAccount);
  }

  function removeNode(address nodeAccount) public onlyOwner {
    for (uint256 i = 0; i <= nodes.length-1; i++){
      if (nodes[i] == nodeAccount) {
        for (uint256 j = i; j < nodes.length-1; j++){
          nodes[j] = nodes[j+1];
        }
        delete nodes[nodes.length-1];
        if (nodes.length == 1) {
          delete nodes;
        }
        break;
      }
    }
  }

  function getNodes() public view returns (address[] memory) {
    return nodes;
  }
}
