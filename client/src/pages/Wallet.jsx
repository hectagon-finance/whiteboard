import React, { useState } from "react";
import { saveAs } from "file-saver";
import elliptic from "elliptic";

const EC = elliptic.ec;
const ec = new EC("p256");

function Wallet() {
  const [privateKey, setPrivateKey] = useState(null);
  const [address, setAddress] = useState(null);
  const [error, setError] = useState(null);

  function handleCreateWallet() {
    const keyPair = ec.genKeyPair();
    const privateKeyHex = keyPair.getPrivate("hex");
    setPrivateKey(privateKeyHex);
    setAddress(null);
    setError(null);
  }

  function handleDownloadPrivateKey() {
    const blob = new Blob([privateKey], { type: "text/plain;charset=utf-8" });
    saveAs(blob, "private-key.txt");
  }

  function handleFileChange(e) {
    const file = e.target.files[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event) => {
        const privateKeyHex = event.target.result.trim();
        if (isValidPrivateKey(privateKeyHex)) {
          const keyPair = ec.keyFromPrivate(privateKeyHex);
          const publicKeyHex = keyPair.getPublic("hex");
          setAddress(generateAddress(publicKeyHex));
          setError(null);
        } else {
          setError("Invalid private key format.");
        }
      };
      reader.readAsText(file);
    }
  }

  function isValidPrivateKey(privateKeyHex) {
    try {
      ec.keyFromPrivate(privateKeyHex);
      return true;
    } catch (err) {
      return false;
    }
  }

  function generateAddress(publicKeyHex) {
    const hash = ec.hash().update(publicKeyHex).digest("hex");
    return "0x" + hash.substring(0, 40);
  }

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <div className="p-4">
        <h1 className="text-2xl font-bold mb-4">Create Wallet</h1>
        {privateKey ? (
          <div className="flex flex-col items-center">
            <p className="mb-4">
              Private key: <code>{privateKey}</code>
            </p>
            <button
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
              onClick={handleDownloadPrivateKey}
            >
              Download
            </button>
            <p className="mt-4">
              Address: <code>{address}</code>
            </p>
            <input
              type="file"
              accept="text/plain"
              onChange={handleFileChange}
              className="mt-4"
            />
            {error && <p className="text-red-500 mt-2">{error}</p>}
          </div>
        ) : (
          <button
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            onClick={handleCreateWallet}
          >
            Create Wallet
          </button>
        )}
      </div>
    </div>
  );
}

export default Wallet;
