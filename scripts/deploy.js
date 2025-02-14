async function main() {
  const SuccorTrail = await ethers.getContractFactory("SuccorTrail");
  const succorTrail = await SuccorTrail.deploy();

  await succorTrail.deployed();

  console.log("SuccorTrail deployed to:", succorTrail.address);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
