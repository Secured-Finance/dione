// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

interface IDioneStaking {
    function mine(address _minerAddr) external;
    function mineAndStake(address _minerAddr) external;
    function stake(uint256 _amount) external;
    function withdraw(uint256 _amount) external;
    function totalStake() external view returns (uint256);
    function minerStake(address _minerAddr) external view returns (uint256);
    function setMinerReward(uint256 _minerReward) external;
    function isMiner(address _minerAddr) external view returns (bool);
    function setMinimumStake(uint256 _minimumStake) external;
    function setAggregator(address _aggregatorAddr) external;
    function slashMiner(address miner, address[] calldata receipentMiners) external;
}
