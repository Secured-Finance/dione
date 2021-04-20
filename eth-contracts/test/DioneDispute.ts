import { ethers } from "hardhat";
import { Signer, Contract } from "ethers";
import { expect } from "chai";
import { soliditySha3 } from "web3-utils";

describe("DioneDispute", function () {
  let dioneDispute: Contract;
  let dioneStaking: Contract;

  beforeEach(async function () {
    const DioneToken = await ethers.getContractFactory("DioneToken");
    const DioneOracle = await ethers.getContractFactory("DioneOracle");
    const DioneDispute = await ethers.getContractFactory("DioneDispute");
    const DioneStaking = await ethers.getContractFactory("DioneStaking");
    const Mediator = await ethers.getContractFactory("Mediator");

    const dioneToken = await DioneToken.deploy();
    await dioneToken.deployed();
    console.log("DioneToken deployed to:", dioneToken.address);

    const _dioneStaking = await DioneStaking.deploy(dioneToken.address, ethers.constants.WeiPerEther.mul(100), 0, ethers.constants.WeiPerEther.mul(5000));
    await _dioneStaking.deployed();
    console.log("staking_contract_address = \"" + _dioneStaking.address+ "\"");

    const _dioneDispute = await DioneDispute.deploy(_dioneStaking.address, 2);
    await _dioneDispute.deployed();
    console.log("dispute_contract_address = \"" + _dioneDispute.address+ "\"");

    const dioneOracle = await DioneOracle.deploy(_dioneStaking.address);
    await dioneOracle.deployed();
    console.log("oracle_contract_address = \"" + dioneOracle.address+ "\"");

    const mediator = await Mediator.deploy(dioneOracle.address);
    await mediator.deployed();
    console.log("mediator_contract_address = \"" + mediator.address +"\"")

    await _dioneStaking.setOracleContractAddress(dioneOracle.address);
    await _dioneStaking.setDisputeContractAddress(_dioneDispute.address);

    const addresses = ["0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "0x70997970c51812dc3a010c7d01b50e0d17dc79c8", "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "0x90F79bf6EB2c4f870365E785982E1f101E93b906"]
    await dioneToken.mint(addresses[0], ethers.constants.WeiPerEther.mul(36000));
    for (const address of addresses) {
      if(address == "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266") {
        continue;
      }
      await dioneToken.transfer(address, ethers.constants.WeiPerEther.mul(9000));
    }

    await dioneToken.transferOwnership(_dioneStaking.address);

    const signers = await ethers.getSigners();
    for (var i = 0; i < addresses.length; i++) {
      const staking = _dioneStaking.connect(signers[i]);
      const token = dioneToken.connect(signers[i]);
      await token.approve(_dioneStaking.address, ethers.constants.WeiPerEther.mul(9000));
      await staking.stake(ethers.constants.WeiPerEther.mul(9000));
      const stake = await _dioneStaking.minerStake(addresses[i]);
      console.log(addresses[i], stake.toString());
    }

    dioneDispute = _dioneDispute;
    dioneStaking = _dioneStaking;
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

    await delay(2000);

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

function delay(ms: number): Promise<void> {
    return new Promise( resolve => setTimeout(resolve, ms) );
}