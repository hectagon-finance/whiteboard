import React, { useState, useRef } from "react";
import { saveAs } from "file-saver";
import elliptic from "elliptic";
import hash from "hash.js";
import { useNavigate } from "react-router-dom";

const EC = elliptic.ec;
const ec = new EC("p256");

function Wallet() {
  const [privateKey, setPrivateKey] = useState(null);
  const [publicKey, setPublicKey] = useState(null);
  const [walletAddress, setWalletAddress] = useState(null);
  const [showInfo, setShowInfo] = useState(false);
  const [validPrivateKey, setValidPrivateKey] = useState(false);
  const [error, setError] = useState(null);
  const fileInputRef = useRef();

  const navigate = useNavigate();

  function handleNavigate() {
    if (validPrivateKey) {
      navigate("/todo-app", {
        state: {
          privateKey: privateKey,
          publicKey: publicKey,
          walletAddress: walletAddress,
        },
      });
    }
  }

  function handleCreateWallet() {
    const keyPair = ec.genKeyPair();
    const privateKeyHex = keyPair.getPrivate("hex");
    const publicKeyHex = keyPair.getPublic("hex").substring(2);
    const generatedWalletAddress = generateWalletAddress(publicKeyHex);
    setPrivateKey(privateKeyHex);
    setWalletAddress(generatedWalletAddress);
    setShowInfo(true);
  }

  function handleDownloadPrivateKey() {
    const blob = new Blob([privateKey], { type: "text/plain;charset=utf-8" });
    saveAs(blob, "private-key.txt");
  }

  function isValidPrivateKey(privateKey) {
    const privateKeyRegex = /^[a-fA-F0-9]{64}$/;
    return privateKeyRegex.test(privateKey);
  }

  function generateWalletAddress(publicKey) {
    const publicKeyHash = hash
      .sha256()
      .update(publicKey)
      .digest("hex")
      .substring(30);
    const walletAddress = "whiteboard" + publicKeyHash;
    return walletAddress;
  }

  async function handleUploadPrivateKey(e) {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = async (e) => {
        const privateKey = e.target.result.trim();
        if (isValidPrivateKey(privateKey)) {
          setError(null);
          setPrivateKey(privateKey);
          const keyPair = ec.keyFromPrivate(privateKey);
          const publicKeyHex = keyPair.getPublic("hex").substring(2);
          setPublicKey(publicKeyHex);
          const generatedWalletAddress = generateWalletAddress(publicKeyHex);
          setWalletAddress(generatedWalletAddress);
          setValidPrivateKey(true);
        } else {
          setError("Invalid private key format.");
        }
      };
      reader.readAsText(file);
    }
  }

  console.log("privateKey", privateKey);
  console.log("publicKey", publicKey);
  console.log("walletAddress", walletAddress);

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        <h1 className="text-2xl font-bold mb-4 text-center">Create Wallet</h1>
        {showInfo ? (
          <div>
            <p className="mb-4 truncate">
              Private key: <code>{privateKey}</code>
            </p>
            <p className="mb-4 truncate">
              Wallet address: <code>{walletAddress}</code>
            </p>
            <button
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded w-full"
              onClick={handleDownloadPrivateKey}
            >
              Download
            </button>
          </div>
        ) : (
          <button
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded w-full"
            onClick={handleCreateWallet}
          >
            Create Wallet
          </button>
        )}
        <input
          ref={fileInputRef}
          type="file"
          accept=".txt"
          className="hidden"
          onChange={handleUploadPrivateKey}
        />
        <div className="flex items-center justify-between mt-4">
          <button
            className={`bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded ${
              validPrivateKey ? "bg-green-500 hover:bg-green-700" : ""
            }`}
            onClick={() => fileInputRef.current.click()}
          >
            {validPrivateKey ? "Logged" : "Login"}
          </button>
          <button
            className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded ml-4"
            onClick={handleNavigate}
            disabled={!validPrivateKey}
          >
            Go to Todo App
          </button>
        </div>
        {error && <p className="text-red-500 mt-2">{error}</p>}
      </div>
    </div>
  );
}

export default Wallet;
