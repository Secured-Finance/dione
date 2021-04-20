import { run, ethers } from "hardhat";
import "@nomiclabs/hardhat-ethers";

async function main() {
    await run("compile");
  
    const DioneToken = await ethers.getContractFactory("DioneToken");
    const DioneOracle = await ethers.getContractFactory("DioneOracle");
    const DioneDispute = await ethers.getContractFactory("DioneDispute");
    const DioneStaking = await ethers.getContractFactory("DioneStaking");
    const Mediator = await ethers.getContractFactory("Mediator");

    const dioneToken = await DioneToken.deploy();
    await dioneToken.deployed();
    console.log("DioneToken deployed to:", dioneToken.address);

    const dioneStaking = await DioneStaking.deploy(dioneToken.address, ethers.constants.WeiPerEther.mul(100), 0, ethers.constants.WeiPerEther.mul(5000));
    await dioneStaking.deployed();
    console.log("staking_contract_address = \"" + dioneStaking.address+ "\"");

    const dioneDispute = await DioneDispute.deploy(dioneStaking.address, 5);
    await dioneDispute.deployed();
    console.log("dispute_contract_address = \"" + dioneDispute.address+ "\"");

    const dioneOracle = await DioneOracle.deploy(dioneStaking.address);
    await dioneOracle.deployed();
    console.log("oracle_contract_address = \"" + dioneOracle.address+ "\"");

    const mediator = await Mediator.deploy(dioneOracle.address);
    await mediator.deployed();
    console.log("mediator_contract_address = \"" + mediator.address +"\"")

    await dioneStaking.setOracleContractAddress(dioneOracle.address);
    await dioneStaking.setDisputeContractAddress(dioneDispute.address);

    const addresses = ["0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266", "0x70997970c51812dc3a010c7d01b50e0d17dc79c8", "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC", "0x90F79bf6EB2c4f870365E785982E1f101E93b906"]
    await dioneToken.mint(addresses[0], ethers.constants.WeiPerEther.mul(50000));
    for (const address of addresses) {
      if(address == "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266") {
        continue;
      }
      await dioneToken.transfer(address, ethers.constants.WeiPerEther.mul(6000));
    }

    await dioneToken.transferOwnership(dioneStaking.address);

    const signers = await ethers.getSigners();
    for (var i = 0; i < addresses.length; i++) {
      const staking = dioneStaking.connect(signers[i]);
      const token = dioneToken.connect(signers[i]);
      await token.approve(dioneStaking.address, ethers.constants.WeiPerEther.mul(5000));
      await staking.stake(ethers.constants.WeiPerEther.mul(5000));
      const stake = await dioneStaking.minerStake(addresses[i]);
      console.log(addresses[i], stake.toString());
    }
  }
  
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });