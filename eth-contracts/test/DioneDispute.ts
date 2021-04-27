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
      nodeCount: 4,
      logging: false,
      minStakeForDisputeVotes: 100
    });

    dioneDispute = contracts.dioneDispute;
    dioneStaking = contracts.dioneStaking;
  });

  it("should create dispute, vote it by various eth addresses and then finish", async function () {
    const addresses = await (await ethers.getSigners()).slice(0, 4);
   
    const dhash = soliditySha3(addresses[1].address, 1);

    await expect(dioneDispute.beginDispute(addresses[1].address, 1))
      .to.emit(dioneDispute, 'NewDispute')
      .withArgs(dhash, 1, addresses[1].address, addresses[0].address);

    for (const x of addresses) {
      if (x == addresses[1] || x == addresses[0]) continue;
      await expect(dioneDispute.connect(x).vote(dhash, true))
        .to.emit(dioneDispute, 'NewVote')
        .withArgs(dhash, x.address);
    }

    await ethers.provider.send("evm_increaseTime", [2]);

    await expect(dioneDispute.finishDispute(dhash))
      .to.emit(dioneDispute, 'DisputeFinished')
      .withArgs(dhash, true);

    expect(await dioneStaking.minerStake(addresses[1].address))
      .to.equal(0);
    
    for (const x of addresses) {
      if (x == addresses[1]) continue;
      expect(await dioneStaking.minerStake(x.address))
        .to.equal(ethers.constants.WeiPerEther.mul(12000));
    }
  });

  it("should fail voting on non-existing dispute", async function() {
    const [addr1] = await ethers.getSigners();
   
    const dhash = soliditySha3(addr1.address, 1);

    await expect(dioneDispute.vote(dhash, true))
      .to.be.revertedWith("dispute doesn't exist");
  });

  it("should fail finishing non-existing dispute", async function() {
    const [addr1] = await ethers.getSigners();
   
    const dhash = soliditySha3(addr1.address, 1);

    await expect(dioneDispute.finishDispute(dhash))
      .to.be.revertedWith("dispute doesn't exist");
  })

  it("should finish dispute with \"false\" result", async function () {
    const addresses = await (await ethers.getSigners()).slice(0, 4);
   
    const dhash = soliditySha3(addresses[1].address, 1);
    await dioneDispute.beginDispute(addresses[1].address, 1);
    await dioneDispute.connect(addresses[2]).vote(dhash, false);
    await dioneDispute.connect(addresses[3]).vote(dhash, false);
    await ethers.provider.send("evm_increaseTime", [2]);
    await expect(dioneDispute.finishDispute(dhash))
      .to.emit(dioneDispute, 'DisputeFinished')
      .withArgs(dhash, false);

    // check if stakes of miners is same as initial
    for (const x of addresses) {
      expect(await dioneStaking.minerStake(x.address))
        .to.be.equal(ethers.constants.WeiPerEther.mul(9000));
    }
  });

  it("should fail when finishing dispute before exceeding vote window time", async () => {
    const [, addr1 ] = await ethers.getSigners();

    const dhash = soliditySha3(addr1.address, 1);
    await dioneDispute.beginDispute(addr1.address, 1);
    await expect(dioneDispute.finishDispute(dhash))
      .to.be.revertedWith("vote window hasn't passed yet");
  });

  it("should fail when voting as dispute initiator", async () => {
    const [, addr1 ] = await ethers.getSigners();

    const dhash = soliditySha3(addr1.address, 1);
    await dioneDispute.beginDispute(addr1.address, 1);
    await expect(dioneDispute.vote(dhash, true))
      .to.be.revertedWith("dispute initiator isn't allowed to vote");
  });

  it("should fail when voting as miner against whom dispute has beginned", async () => {
    const [, addr1 ] = await ethers.getSigners();

    const dhash = soliditySha3(addr1.address, 1);
    await dioneDispute.beginDispute(addr1.address, 1);
    await expect(dioneDispute.connect(addr1).vote(dhash, true))
      .to.be.revertedWith("the miner against whom dispute has beginned isn't allowed to vote");
  });

  describe("DioneDispute - insufficient funds", () => {
    let dioneDispute: Contract;
    before(async () => {
      const contracts = await deploy({
        reward: 100,
        minStake: 0,
        voteWindowTime: 2,
        randomizeStake: false,
        maxStake: 0, // don't use this deployment feature
        actualStake: 50,
        nodeCount: 1,
        logging: false,
        minStakeForDisputeVotes: 100
      });
  
      dioneDispute = contracts.dioneDispute;
    });

    it("should fail when voting with insufficient funds", async () => {
      const [, addr1, addr2] = await ethers.getSigners();

      const dhash = soliditySha3(addr1.address, 1);
      await dioneDispute.beginDispute(addr1.address, 1);
      await expect(dioneDispute.connect(addr2).vote(dhash, true))
        .to.be.revertedWith("miner doesn't have minimum stake to vote");
    });
  });
});