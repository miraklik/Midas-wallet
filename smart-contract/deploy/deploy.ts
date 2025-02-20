import { Wallet } from "zksync-ethers";
import { HardhatRuntimeEnvironment } from "hardhat/types";
import { Deployer } from "@matterlabs/hardhat-zksync";
import { vars } from "hardhat/config";

export default async function (hre: HardhatRuntimeEnvironment) {
  console.log(`Running deploy script`);

  const wallet = new Wallet(vars.get("PRIVATE_KEY"));

  const deployer = new Deployer(hre, wallet);
  const artifact = await deployer.loadArtifact("Wallet");

  const tokenContract = await deployer.deploy(artifact);

  console.log(
    `${
      artifact.contractName
    } was deployed to ${await tokenContract.getAddress()}`
  );
}