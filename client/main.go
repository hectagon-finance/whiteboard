package main

import (
	// "bytes"
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/hectagon-finance/whiteboard/crypto"
	"github.com/hectagon-finance/whiteboard/types"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	// gob.Register(elliptic.P256())
	if os.Args[1] == "create-wallet" {
		//check client have wallet => can't create
		if checkHaveWallet("public_key.txt", true) == false {
			privateKey := crypto.GeneratePrivateKey()
			publicKey := privateKey.PublicKey()
			address := publicKey.Address()
			privateKeyStr := privateKey.PrivateKeyStr()
			publicKeyStr := (publicKey.PublicKeyStr())
			fmt.Println("Please save your private key:")
			fmt.Println("Private Key:", privateKeyStr)
			fmt.Println("Public Key:", publicKeyStr)
			fmt.Println("Address is:", address)

			// create file txt in root folder
			file, err := os.Create("public_key.txt")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			// write publicKeyStr to file
			_, err = file.WriteString(publicKeyStr)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("public_key file created successfully.")
		}
	}

	if os.Args[1] == "send" && os.Args[3] == "-k" {
		if checkHaveWallet("public_key.txt", false) {
			// Open file public_key.txt
			file, err := os.Open("public_key.txt")
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			// Read data from file and assign for publicKeyStr
			scanner := bufio.NewScanner(file)
			scanner.Scan()
			publicKeyStr := scanner.Text()

			publicKey := crypto.PublicKeyFromString(publicKeyStr)
			publicKeyForConvert := crypto.PublicKeyFromString(publicKeyStr).Key
			privateKey := crypto.PrivateKeyFromString(os.Args[4], publicKeyForConvert)

			msg := []byte(os.Args[2])

			sig, err := privateKey.Sign(msg)

			// input := Input{
			// 	PublicKey: publicKey,
			// 	SigNature: sig,
			// 	Data:      msg,
			// }

			// var buf bytes.Buffer
			// enc := gob.NewEncoder(&buf)

			// if err := enc.Encode(input); err != nil {
			// 	return
			// }

			flag.Parse()
			log.SetFlags(0)

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt)
			
			haha := *publicKey
			hehe := *sig
			tx := types.NewTransaction(haha, hehe, msg)
			sendTransaction("9000", tx)

		}
	}
}

func checkHaveWallet(s string, showMessage bool) bool {
	_, err := os.Stat(s)
	if err == nil {
		if showMessage {
			fmt.Printf("You have already created a wallet, if you want to recreate it, please delete the file %s in the original directory \n", s)
			return true
		} else {
			return true
		}
	} else {
		if showMessage {
			return false
		} else {
			fmt.Println("Something went wrong! if you, if you haven't created a wallet please create one")
			return false
		}
	}
}

func sendTransaction(validatorId string, tx types.Transaction) {
	u := url.URL{Scheme: "ws", Host: "localhost:" + validatorId, Path: "/ws"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Println("Error connecting to the validator:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Sending transaction to validator", validatorId)

	publicKey := tx.PublicKey()
	publicKeyStr := publicKey.PublicKeyStr()

	signature := tx.Signature()
	signatureStr := signature.SignatureStr()

	message := map[string]interface{}{
		"type":          "transaction",
		"from":          "client",
		"validatorId":   "fake",
		"transactionId": tx.Id(),
		"timestamp":     tx.Timestamp(),
		"publicKey" :    publicKeyStr,
		"signature" :    signatureStr,
		"data":          string(tx.Data()),
	}
	
	msg, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling the message:", err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, msg)
}