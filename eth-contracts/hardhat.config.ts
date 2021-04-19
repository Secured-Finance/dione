import { task } from "hardhat/config";
import "@nomiclabs/hardhat-ethers";

task("accounts", "Prints the list of accounts", async (args, hre) => {
  const accounts = await hre.ethers.getSigners();

  for (const account of accounts) {
    console.log(await account.address);
  }
});

export default {
  solidity: "0.8.3",
  networks: {
    ganache: {
      url: `http://localhost:7545`,
      accounts: {
        mnemonic: "test test test test test test test test test test test junk"
      }
    }
  }
};