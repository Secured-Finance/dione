<h1 align="center">
  <img  src="assets/img/dione.jpg" alt="dione" />
</h1>

<h3 align="center">Secured Finance - A protocol for financial transactions backed by crypto-assets.</h3>

[![badge](https://img.shields.io/badge/submit%20for-HackFS-blue)](https://hack.ethglobal.co/showcase/secured-finance-recTkx6c1RDoLeaQm)

# Table of Contents

- [Introduction](#introduction)
- [Specifications](#specifications)
  - [State validation](#state-validation)
  - [Connectors](#connectors)
  - [External adapters](#external-adapters)
  - [Consensus](#consensus)
- [DIONE token](#dione-token)
- [Contribute](#contribute)
- [License](#license)


# Introduction

Dione oracle network is solving the interoperability problem between multiple blockchain networks. Blockchains themself could be describled as deterministic state machine replicated on full-nodes that retains consensus safety. And having the network that handles state changes validation across multiple protocols is crucial for [Secured Finance](https://github.com/Secured-Finance) first-layer transactions execution. Using Dione network Secured Finance could provide completely decentralized service for it's users.

Simple explanation for Dione is a blockchain that tracks the state of another blockchains. Dione network itself has fast-finality because of the pBFT (Practical Byzantine Fault Tolerance) consensus and has a role to establish finality for the blockchain it connects (Ethereum, Bitcoin networks has probabilistic-finality). 

The networking layer of Dione node is based on [libp2p](https://github.com/libp2p/go-libp2p), a flexible cross-platform network framework for peer-to-peer applications. Projects like [ETH2 (Ethereum 2.0)](https://github.com/ethereum/eth2.0-specs), [Filecoin](https://github.com/filecoin-project) and Polkadot's [Substrate](https://github.com/paritytech/substrate) is based on libp2p making it's the standard for future decentralized infrastructures.

# Specifications
## State validation

In order to validate the state from another blockchain in the Ethereum smart contract Dione users has to import [OracleEmitter.sol](https://github.com/Secured-Finance/p2p-oracle-smart-contracts/blob/master/contracts/OracleEmitter.sol) in smart contract that has to access external data. 

## Connectors

Most of connections established via RPC request–response protocol meaning that Dione network has to meet all interface requirements and be customized for the particular chain it's connects to.

## External adapters

Most of off-chain data use cases require external adapters and APIs for simple integration of custom data sources. However Dione network aiming to solve cross-chain interoperability external adapters would be implemented in the second phase using REST-based communication.

## Consensus

Dione network uses PBFT a practical Byzantine fault-tolerant consensus protocol invented by Miguel Castro and Barbara Liskov at MIT. At a high level, it operates by running a leader election in every block in which, on expectation, a set number of participants may be eligible to submit a block. 

The randomness used in the proofs is generated from [DRAND](https://drand.love), an unbiasable randomness generator, through a beacon.

In order to identify which block has to be selected between multiple blocks proposed by multiple leaders other network operators has to vote up to 2/3 of total amount of nodes running in the network. This enforces Dione network as fast-finality chain.

# DIONE token
DIONE token is a ERC20 token on Ethereum blockchain used as a staking token to run the Dione node and for decentralized governance. 

Staking mechanism would encourage node operators to behave correctly as well as malicious nodes would lose their stakes. In order to run Dione node the operator has to stake at least 10,000 DIONE tokens. 

Total supply of DIONE tokens is 2,000,000. By that the maximum amount of nodes participating in Dione is limited to 200 at the initial start of Dione network. 30% of DIONE tokens would be allocated to Secured Finance team.

Governance mechanism would take place on [snapshot.page](https://snapshot.page/#/). The main application of governance descision is around required amount of DIONE tokens to stake in order to run Dione node. Additional governance descision could be proposed by DIONE token holders among the way.


# Contribute <a name="contribute"> </a> 

#### Dione core contributors:
[Denis Davydov](https://github.com/ChronosX88)

[Bach Adylbekov](https://github.com/bahadylbekov)

We welcome every contributions big and small! Take a look at the [community contributing notes](). Please make sure to check the [issues](https://github.com/Secured-Finance/dione/issues). Search the closed ones before reporting things, and help us with the open ones.


# License

This project is licensed under the MIT license, Copyright (c) 2020 Secured Finance. For more information see `LICENSE.md`.
