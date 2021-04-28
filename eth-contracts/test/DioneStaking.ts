import { BigNumber } from "@ethersproject/bignumber";
import { Contract } from "@ethersproject/contracts";
import { SignerWithAddress } from "@nomiclabs/hardhat-ethers/signers";
import { expect } from "chai";
import { ethers } from "hardhat";

describe("DioneStaking", () => {
    let signers: SignerWithAddress[];
    let allSigners: SignerWithAddress[];
    let dioneToken: Contract;
    let dioneStaking: Contract;
    const rewardAmount = 100;
    const minStakeAmount = 500;

    beforeEach(async () => {
        allSigners = await ethers.getSigners();
        signers = allSigners.slice(0, 4);

        const DioneToken = await ethers.getContractFactory("DioneToken");
        const DioneStaking = await ethers.getContractFactory("DioneStaking");

        dioneToken = await DioneToken.deploy();
        await dioneToken.deployed();

        dioneStaking = await DioneStaking.deploy(
            dioneToken.address, 
            toWei(rewardAmount), 
            0, 
            toWei(minStakeAmount)
        );
        await dioneStaking.deployed();

        for (var s of signers) {
            await dioneToken.mint(s.address, toWei(5000));
        }

        await dioneToken.transferOwnership(dioneStaking.address);
    });

    it("should successfully stake some of tokens", async () => {
        for (var s of signers) {
            await dioneToken.connect(s).approve(dioneStaking.address, toWei(666));
            await expect(dioneStaking.connect(s).stake(toWei(666)))
                .to.emit(dioneStaking, "Stake")
                .withArgs(s.address, toWei(666));
            expect(await dioneStaking.minerStake(s.address))
                .to.be.equal(toWei(666));
            expect(await dioneStaking.isMiner(s.address))
                .to.be.equal(true);
        }
    });
    
    it("check miner balance without stake", async () => {
        expect(await dioneStaking.minerStake(signers[0].address))
            .to.be.equal(toWei(0));
        expect(await dioneStaking.isMiner(signers[0].address))
            .to.be.equal(false);
    });

    it("should mine reward successfully", async () => {
        await dioneStaking.setOracleContractAddress(signers[0].address)

        await dioneToken.approve(dioneStaking.address, toWei(500));
        await dioneStaking.stake(toWei(500));
        const res = await dioneStaking.mine(signers[0].address);
        await expect(res)
            .to.emit(dioneStaking, "Mine")
            .withArgs(signers[0].address, res.blockNumber);

        expect(await dioneToken.balanceOf(signers[0].address))
            .to.be.equal(toWei(4600));
    });

    it("should mine and stake reward successfully", async () => {
        await dioneStaking.setOracleContractAddress(signers[0].address);

        await dioneToken.approve(dioneStaking.address, toWei(500));
        await dioneStaking.stake(toWei(500));

        await dioneStaking.mineAndStake(signers[0].address);
        expect(await dioneStaking.minerStake(signers[0].address))
            .to.be.equal(toWei(600));
    });

    it("should fail if caller of mine methods isn't desired address", async () => {
        await expect(dioneStaking.mine(signers[0].address))
            .to.be.revertedWith("not oracle contract");
        await expect(dioneStaking.mineAndStake(signers[0].address))
            .to.be.revertedWith("not oracle contract");
    });

    it("check change of total stake after staking", async () => {
        const expectedTotalStake = toWei(666*4);

        for (var s of signers) {
            await dioneToken.connect(s).approve(dioneStaking.address, toWei(666));
            await dioneStaking.connect(s).stake(toWei(666));
        }

        expect(await dioneStaking.totalStake())
            .to.be.equal(expectedTotalStake);
    });

    it("should withdraw funds successfully", async () => {
        for (var s of signers) {
            await dioneToken.connect(s).approve(dioneStaking.address, toWei(500));
            await dioneStaking.connect(s).stake(toWei(500));
        }

        for (var s of signers) {
            await expect(dioneStaking.connect(s).withdraw(toWei(500)))
                .to.emit(dioneStaking, "Withdraw")
                .withArgs(s.address, toWei(500));

            expect(await dioneToken.balanceOf(s.address))
                .to.be.equal(toWei(5000));
        }
    });

    it("should fail withdrawing invalid values", async () => {
        await dioneToken.approve(dioneStaking.address, toWei(500));
        await dioneStaking.stake(toWei(500));

        await expect(dioneStaking.withdraw(toWei(666)))
            .to.be.revertedWith("withdraw: not enough tokens");

        await expect(dioneStaking.withdraw(0))
            .to.be.revertedWith("cannot withdraw zero");
    });

    it("check changing miner reward", async () => {
        await expect(dioneStaking.setMinerReward(0))
            .to.be.revertedWith("reward must not be zero");
        
        await expect(dioneStaking.setMinerReward(toWei(5)))
            .to.emit(dioneStaking, "RewardChanged")
            .withArgs(toWei(100), toWei(5));

        expect(await dioneStaking.minerReward())
            .to.be.equal(toWei(5));

        await expect(dioneStaking.connect(signers[1]).setMinerReward(toWei(5)))
            .to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("check changing minimum stake", async () => {
        await expect(dioneStaking.setMinimumStake(0))
            .to.be.revertedWith("minimum stake must not be zero");
        
        await expect(dioneStaking.setMinimumStake(toWei(5)))
            .to.emit(dioneStaking, "MinimumStakeChanged")
            .withArgs(toWei(500), toWei(5));

        expect(await dioneStaking.minimumStake())
            .to.be.equal(toWei(5));

        await expect(dioneStaking.connect(signers[1]).setMinimumStake(toWei(5)))
            .to.be.revertedWith("Ownable: caller is not the owner");
    });

    it("should fail when staking invalid values", async () => {
        await expect(dioneStaking.stake(0))
            .to.be.revertedWith("cannot stake zero");

        await expect(dioneStaking.stake(toWei(400)))
            .to.be.revertedWith("actual stake amount is less than minimum stake amount");
    });

    it("should fail when calling slashMiner with invalid values (initial checks)", async () => {
        const zeroAddress = "0x0000000000000000000000000000000000000000";

        await expect(dioneStaking.slashMiner(zeroAddress, [signers[0].address]))
            .to.be.revertedWith("caller is not the dispute contract");

        await dioneStaking.setDisputeContractAddress(signers[0].address);

        await expect(dioneStaking.slashMiner(zeroAddress, [signers[0].address]))
            .to.be.revertedWith("slashing address must not be zero");

        await expect(dioneStaking.slashMiner(signers[1].address, [signers[0].address]))
            .to.be.revertedWith("slashing address isn't dione miner");

        for (var s of signers) {
            await dioneToken.connect(s).approve(dioneStaking.address, toWei(500));
            await dioneStaking.connect(s).stake(toWei(500));
        }

        await expect(dioneStaking.slashMiner(signers[1].address, [signers[0].address, signers[2].address, allSigners[4].address]))
            .to.be.revertedWith("receipent address isn't dione miner");
        
        await expect(dioneStaking.slashMiner(signers[1].address, [signers[0].address, signers[2].address, signers[1].address]))
            .to.be.revertedWith("receipent address must not be slashing address");
    });

    it("should fail set*Address methods when passing invalid values", async () => {
        const zeroAddress = "0x0000000000000000000000000000000000000000";

        await expect(dioneStaking.setDisputeContractAddress(zeroAddress))
            .to.be.revertedWith("address must not be zero");
        await expect(dioneStaking.setOracleContractAddress(zeroAddress))
            .to.be.revertedWith("address must not be zero");
    });

    it("should stake few time successfully", async () => {
        const expectedStake = toWei(500*3);

        for (var i = 0; i < 3; i++) {
            await dioneToken.approve(dioneStaking.address, toWei(500));
            await dioneStaking.stake(toWei(500));
        }

        expect(await dioneStaking.minerStake(signers[0].address))
            .to.be.equal(expectedStake);
    });
})

function toWei(eth: number): BigNumber {
    return ethers.constants.WeiPerEther.mul(eth);
}