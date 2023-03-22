package main

import (
	"github.com/hectagon-finance/whiteboard/types"
	"github.com/hectagon-finance/whiteboard/utils/crypto"
)

func main() {
	privateKey := crypto.GeneratePrivateKey()
	publicKey := privateKey.PublicKey()

	msg := []byte("Hello World")

	sig, err := privateKey.Sign(msg)
	if err != nil {
		panic(err)
	}

	trans := types.NewTransaction(publicKey, *sig, msg)

	transactions := []types.Transaction{trans}

	bc := types.NewBlockchain()
	prevHash := bc.LastBlock().GetHash()
	bc.CreateBlock(1, prevHash, transactions)

	prevHash = bc.LastBlock().GetHash()
	bc.CreateBlock(2, prevHash, transactions)

	prevHash = bc.LastBlock().GetHash()
	bc.CreateBlock(3, prevHash, transactions)
	bc.Print()
}
