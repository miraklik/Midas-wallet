import { ethers } from 'ethers';

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "generateWallet") {
    const wallet = ethers.Wallet.createRandom();
    sendResponse({
      address: wallet.address,
      privateKey: wallet.privateKey
    });
  }
});

await chrome.storage.local.set({
    wallet: {
      address: "0x...",
      privateKey: "encrypted_private_key"
    }
  });
  
  const data = await chrome.storage.local.get("wallet");