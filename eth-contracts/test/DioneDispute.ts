import { ethers } from "hardhat";
import { Contract } from "ethers";
import { expect } from "chai";
import { soliditySha3 } from "web3-utils";
import deploy from "../common/deployment";

describe("DioneDispute", function () {
  let dioneDispute: Contract;
  let dioneStaking: Contract;

  beforeEach(async function () {
    const contracts = await deploy({
      reward: 100,
      minStake: 5000,
      voteWindowTime: 2,
      randomizeStake: false,
      maxStake: 0, // don't use this deployment feature
      actualStake: 9000,
      nodeCount: 4
    });

    dioneDispute = contracts.dioneDispute;
    dioneStaking = contracts.dioneStaking;
  });

  it("should create dispute, vote it by various eth addresses and then finish it after 2 secs", async function () {
    const [owner, addr1, addr2, addr3] = await ethers.getSigners();
   
    const dhash = soliditySha3(addr1.address, 1);

    await expect(dioneDispute.beginDispute(addr1.address, 1))
      .to.emit(dioneDispute, 'NewDispute')
      .withArgs(dhash, 1, addr1.address, owner.address);

    await expect(dioneDispute.connect(addr2).vote(dhash, true))
      .to.emit(dioneDispute, 'NewVote')
      .withArgs(dhash, addr2.address);

    await expect(dioneDispute.connect(addr3).vote(dhash, true))
      .to.emit(dioneDispute, 'NewVote')
      .withArgs(dhash, addr3.address);

    await ethers.provider.send("evm_increaseTime", [2]);

    await expect(dioneDispute.finishDispute(dhash))
      .to.emit(dioneDispute, 'DisputeFinished')
      .withArgs(dhash, true);

    expect(await dioneStaking.minerStake(addr1.address))
      .to.equal(0);

    expect(await dioneStaking.minerStake(addr2.address))
      .to.equal(ethers.constants.WeiPerEther.mul(12000));
      
    expect(await dioneStaking.minerStake(addr3.address))
      .to.equal(ethers.constants.WeiPerEther.mul(12000));

    expect(await dioneStaking.minerStake(owner.address))
      .to.equal(ethers.constants.WeiPerEther.mul(12000));
  });
});