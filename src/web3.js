import Web3 from 'web3';
import SuccorTrailABI from './contracts/SuccorTrail.json'; // Assuming this contains JUST the ABI

const web3 = new Web3(Web3.givenProvider || 'https://rpc.testnet.lisk.com'); // Lisk Testnet RPC

// Manually specify contract address (from Remix deployment)
const contractAddress = "0x6b5bf20d25b719252933ee7a09555c9050bc7aee"; 

// Create contract instance
const contract = new web3.eth.Contract(
  SuccorTrailABI.abi, // or SuccorTrailABI if you exported the ABI array directly
  contractAddress
);

export default contract;