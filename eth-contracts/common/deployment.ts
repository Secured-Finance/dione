import { Contract } from "@ethersproject/contracts";
import { ethers } from "hardhat";

interface Environment {
    dioneToken: Contract;
    dioneStaking: Contract;
    dioneDispute: Contract;
    dioneOracle: Contract;
    mediator: Contract;
}

interface DeploymentOptions {
    reward: number;
    minStake: number;
    voteWindowTime: number; // in seconds
    randomizeStake: boolean;
    maxStake: number; // for randomizer
    actualStake: number; // of each node
    nodeCount: number;
    logging: boolean;
}

async function deploy(opts: DeploymentOptions): Promise<Environment> {
    const logger = new LogWrapper(opts.logging);

    const accounts = (await ethers.getSigners()).slice(0, opts.nodeCount);

    const DioneToken = await ethers.getContractFactory("DioneToken");
    const DioneOracle = await ethers.getContractFactory("DioneOracle");
    const DioneDispute = await ethers.getContractFactory("DioneDispute");
    const DioneStaking = await ethers.getContractFactory("DioneStaking");
    const Mediator = await ethers.getContractFactory("Mediator");

    const dioneToken = await DioneToken.deploy();
    await dioneToken.deployed();
    logger.log("DioneToken deployed to:", dioneToken.address);

    const dioneStaking = await DioneStaking.deploy(dioneToken.address, ethers.constants.WeiPerEther.mul(opts.reward), 0, ethers.constants.WeiPerEther.mul(opts.minStake));
    await dioneStaking.deployed();
    logger.log("staking_contract_address = \"" + dioneStaking.address+ "\"");

    const dioneDispute = await DioneDispute.deploy(dioneStaking.address, opts.voteWindowTime);
    await dioneDispute.deployed();
    logger.log("dispute_contract_address = \"" + dioneDispute.address+ "\"");

    const dioneOracle = await DioneOracle.deploy(dioneStaking.address);
    await dioneOracle.deployed();
    logger.log("oracle_contract_address = \"" + dioneOracle.address+ "\"");

    const mediator = await Mediator.deploy(dioneOracle.address);
    await mediator.deployed();
    logger.log("mediator_contract_address = \"" + mediator.address +"\"")

    const env: Environment = {
        dioneToken: dioneToken,
        dioneStaking: dioneStaking,
        dioneDispute: dioneDispute,
        dioneOracle: dioneOracle,
        mediator: mediator
    }

    await dioneStaking.setOracleContractAddress(dioneOracle.address);
    await dioneStaking.setDisputeContractAddress(dioneDispute.address);

    const stakeForEach: number[] = []
    var mintValue = opts.actualStake*opts.nodeCount
    if (opts.randomizeStake) {
        var sum = 0;
        for (var i = 0; i < opts.nodeCount; i++) {
            stakeForEach.push(randomInt(opts.minStake, opts.maxStake));
            sum += stakeForEach[i];
        }
        mintValue = sum;
    }

    await dioneToken.mint(accounts[0].address, ethers.constants.WeiPerEther.mul(mintValue));
    
    var stakeValue = opts.actualStake;
    for (var i = 0; i < accounts.length; i++) {
        if(accounts[i] == accounts[0]) {
            continue;
        }

        if (opts.randomizeStake) {
            stakeValue = stakeForEach[i];
        }

        await dioneToken.transfer(accounts[i].address, ethers.constants.WeiPerEther.mul(stakeValue));
    }

    await dioneToken.transferOwnership(dioneStaking.address);

    for (var i = 0; i < accounts.length; i++) {
        const staking = dioneStaking.connect(accounts[i]);
        const token = dioneToken.connect(accounts[i]);
        if (opts.randomizeStake) {
            stakeValue = stakeForEach[i];
        }
        await token.approve(dioneStaking.address, ethers.constants.WeiPerEther.mul(stakeValue));
        await staking.stake(ethers.constants.WeiPerEther.mul(stakeValue));
        const stake = await dioneStaking.minerStake(accounts[i].address);
        logger.log(accounts[i].address, stake.toString());
    }

    return env;
}

export default deploy;

// min and max included
function randomInt(min: number, max: number){
    return Math.floor(Math.random() * (max - min + 1)) + min;
}

class LogWrapper {
    private enabled: boolean;

    constructor(_enabled: boolean) {
        this.enabled = _enabled;
    }

    public log(message?: any, ...optionalParams: any[]): void {
        if(this.enabled) console.log(message, optionalParams);
    }
}