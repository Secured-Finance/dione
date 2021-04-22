import { task } from "hardhat/config";
import "@nomiclabs/hardhat-ethers";
import "hardhat-tracer";
import "@nomiclabs/hardhat-waffle";
import "hardhat-ethernal";

task("accounts", "Prints the list of accounts", async (args, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(await account.address);
  }
});

export default {
  solidity: "0.8.3",
  networks: {
    geth: {
      url: `http://localhost:8545`,
      accounts: {
        mnemonic: "test test test test test test test test test test test junk"
      }
    }
  }
};