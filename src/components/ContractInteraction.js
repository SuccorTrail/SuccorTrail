import React, { useState, useEffect } from 'react';
import contract from '../web3';

const ContractInteraction = () => {
  const [message, setMessage] = useState('');
  const [newMessage, setNewMessage] = useState('');

  useEffect(() => {
    const fetchMessage = async () => {
      const message = await contract.methods.message().call();
      setMessage(message);
    };
    fetchMessage();
  }, []);

  const updateMessage = async () => {
    const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
    await contract.methods.setMessage(newMessage).send({ from: accounts[0] });
    setMessage(newMessage);
  };

  return (
    <div>
      <h1>Message: {message}</h1>
      <input
        type="text"
        value={newMessage}
        onChange={(e) => setNewMessage(e.target.value)}
      />
      <button onClick={updateMessage}>Update Message</button>
    </div>
  );
};

export default ContractInteraction;