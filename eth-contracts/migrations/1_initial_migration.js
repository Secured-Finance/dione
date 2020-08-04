const Whitelist = artifacts.require("Whitelist");
const OracleEmitter = artifacts.require("OracleEmitter");
const Aggregator = artifacts.require("Aggregator");
const Mediator = artifacts.require("Mediator");

module.exports = function (deployer) {
  return deployer
    .then(() => {
      return deployer.deploy(Whitelist);
    })
    .then(() => {
      return deployer.deploy(OracleEmitter);
    }).then(() => {
      return deployer.deploy(Aggregator, Whitelist.address);
    }).then(() => {
      return deployer.deploy(Mediator, OracleEmitter.address, Aggregator.address);
    });
};
