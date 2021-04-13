const DioneToken = artifacts.require("DioneToken");
const DioneOracle = artifacts.require("DioneOracle");
const DioneDispute = artifacts.require("DioneDispute");
const DioneStaking = artifacts.require("DioneStaking");

module.exports = function (deployer) {
  return deployer
    .then(() => {
      return deployer.deploy(DioneToken);
    })
    .then(() => {
      return deployer.deploy(DioneStaking, DioneToken.address, web3.utils.toBN("10000000000000000000"), 0, web3.utils.toBN("5000000000000000000000"));
    }).then(() => {
      return deployer.deploy(DioneOracle, DioneStaking.address);
    }).then(() => {
      return deployer.deploy(DioneDispute, DioneStaking.address);
    });
};
