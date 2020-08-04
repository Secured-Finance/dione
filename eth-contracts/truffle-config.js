require('dotenv').config();
const customProvider = (mnemonic, rpcEndpoint) => () =>
  new HDWalletProvider(mnemonic, rpcEndpoint);

const infuraProvider = (network) =>
  customProvider(
    process.env.MNEMONIC || process.env.PRIVATE_KEY || '',
    `https://${network}.infura.io/v3/${process.env.WEB3_INFURA_ID}`,
  );

const ropstenProvider = infuraProvider('ropsten');

module.exports = {
  networks: {
    development: {
      host: "127.0.0.1",     // Localhost (default: none)
      port: 7545,            // Standard Ethereum port (default: none)
      network_id: "*",       // Any network (default: none)
    },
    ropsten: {
      provider: ropstenProvider,
      network_id: 3,
      // gasPrice: 5000000000,
      // gas: 4500000,
      // gasPrice: 10000000000,
      // confirmations: 0, // # of confs to wait between deployments. (default: 0)
      skipDryRun: true,
    },
  },
  compilers: {
    solc: {
      version: '0.6.12',
      settings: {
        optimizer: {
          enabled: true, // Default: false
          runs: 0, // Default: 200
        },
      },
    },
  }
};
