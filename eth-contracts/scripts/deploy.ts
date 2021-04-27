import { run } from "hardhat";
import "@nomiclabs/hardhat-ethers";
import deploy from "../common/deployment";

async function main() {
    await run("compile");
    
    await deploy({
      reward: 100,
      minStake: 5000,
      voteWindowTime: 5,
      randomizeStake: false,
      maxStake: 0, // don't use this deployment feature
      actualStake: 5000,
      nodeCount: 4,
      logging: true,
      minStakeForDisputeVotes: 100
    });
  }
  
  main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });