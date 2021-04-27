import { ethers } from "hardhat";
import { BigNumber, Contract, providers, utils } from "ethers";
import { expect } from "chai";
import deploy from "../common/deployment";

describe("DioneOracle", function () {
  let dioneOracle: Contract;
  let dioneToken: Contract;

  beforeEach(async function () {
    const contracts = await deploy({
      reward: 100,
      minStake: 5000,
      voteWindowTime: 2,
      randomizeStake: false,
      maxStake: 0, // don't use this deployment feature
      actualStake: 9000,
      nodeCount: 4,
      logging: false,
      minStakeForDisputeVotes: 100
    });
    dioneOracle = contracts.dioneOracle;
    dioneToken = contracts.dioneToken;
  });

  it("should create request and cancel it", async function () {
    const timestamp = 1625097600;
    await ethers.provider.send("evm_setNextBlockTimestamp", [timestamp]);
    const requestDeadline = timestamp + 300;
    await expect(dioneOracle.requestOracles(1, "getTransaction", "bafy2bzaceaaab3kkoaocal2dzh3okzy4gscqpdt42hzrov3df6vjumalngc3g", dioneOracle.address, BigNumber.from(0x8da5cb5b)))
      .to.emit(dioneOracle, 'NewOracleRequest')
      .withArgs(1, "getTransaction", "bafy2bzaceaaab3kkoaocal2dzh3okzy4gscqpdt42hzrov3df6vjumalngc3g", 1, requestDeadline);

    await expect(dioneOracle.cancelOracleRequest(1))
      .to.emit(dioneOracle, 'CancelOracleRequest')
      .withArgs(1);

    // let's check whether submission fails
    await expect(dioneOracle.submitOracleRequest(1, BigNumber.from(0x8da5cb5b)))
      .to.be.revertedWith("this request is not pending");
  });

  it("should create request and submit it", async function() {
    const [addr0] = await ethers.getSigners();

    await dioneOracle.requestOracles(1, "getTransaction", "bafy2bzaceaaab3kkoaocal2dzh3okzy4gscqpdt42hzrov3df6vjumalngc3g", dioneOracle.address, BigNumber.from(0x8da5cb5b));
    const res = dioneOracle.submitOracleRequest(1, BigNumber.from(0x8da5cb5b));
    await expect(res)
      .to.emit(dioneOracle, "SubmittedOracleRequest")
      .withArgs(1, BigNumber.from(0x8da5cb5b));

    // check if miner has received the reward
    expect(await dioneToken.balanceOf(addr0.address))
      .to.be.equal(ethers.constants.WeiPerEther.mul(100));
  });

  it("should fail submission after request deadline", async function () {
    await dioneOracle.requestOracles(1, "getTransaction", "bafy2bzaceaaab3kkoaocal2dzh3okzy4gscqpdt42hzrov3df6vjumalngc3g", dioneOracle.address, BigNumber.from(0x8da5cb5b));

    await ethers.provider.send("evm_increaseTime", [301]);
    await expect(dioneOracle.submitOracleRequest(1, BigNumber.from(0x8da5cb5b)))
      .to.be.revertedWith("submission has exceeded the deadline");
  });

  it("should fail submission on invalid request id", async function () {
    await expect(dioneOracle.submitOracleRequest(333, BigNumber.from(0x8da5cb5b)))
      .to.be.revertedWith("this request is not pending");
  });

  it("should fail cancel of request with invalid request id", async function () {
    await expect(dioneOracle.cancelOracleRequest(333))
      .to.be.revertedWith("this request is not pending");
  });
});