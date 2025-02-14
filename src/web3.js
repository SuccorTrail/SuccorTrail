import Web3 from 'web3';
import SuccorTrail from './contracts/SuccorTrail.json';

const web3 = new Web3(Web3.givenProvider || 'http://localhost:8545');
const networkId = await web3.eth.net.getId();
const deployedNetwork = SuccorTrail.networks[networkId];
const contract = new web3.eth.Contract(
  SuccorTrail.abi,
  deployedNetwork && deployedNetwork.address,
);

export default contract;