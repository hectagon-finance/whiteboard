package main

import (
	// "bytes"
	"bufio"
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hectagon-finance/whiteboard/crypto"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

type Input struct {
	PublicKey *crypto.PublicKey
	SigNature *crypto.Signature
	Data      []byte
}

func main() {
	gob.Register(elliptic.P256())
	if os.Args[1] == "create-wallet" {
		//check client have wallet => can't create
		if checkHaveWallet("public_key.txt", true) == false {
			privateKey := crypto.GeneratePrivateKey()
			publicKey := privateKey.PublicKey()
			privateKeyStr := privateKey.PrivateKeyStr()
			publicKeyStr := (publicKey.PublicKeyStr())
			fmt.Println("Please save your private key:")
			fmt.Println("Private Key:", privateKeyStr)
			fmt.Println("Public Key:", publicKeyStr)

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

			msg := []byte(string(os.Args[2]))

			sig, err := privateKey.Sign(msg)

			input := Input{
				PublicKey: publicKey,
				SigNature: sig,
				Data:      msg,
			}

			var buf bytes.Buffer
			enc := gob.NewEncoder(&buf)

			if err := enc.Encode(input); err != nil {
				return
			}

			flag.Parse()
			log.SetFlags(0)

			interrupt := make(chan os.Signal, 1)
			signal.Notify(interrupt, os.Interrupt)

			u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
			log.Printf("connecting to %s", u.String())

			c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Fatal("dial:", err)
			}

			err = c.WriteMessage(websocket.TextMessage, buf.Bytes())
			if err != nil {
				log.Fatal("write:", err)
			}

			defer c.Close()

			done := make(chan struct{})

			go func() {
				defer close(done)
				for {
					_, message, err := c.ReadMessage()
					if err != nil {
						log.Println("read:", err)
						return
					}
					log.Printf("recv: %s", message)
				}
			}()

			for {
				select {
				case <-done:
					return
				case <-interrupt:
					log.Println("interrupt")

					// Cleanly close the connection by sending a close message and then
					// waiting (with timeout) for the server to close the connection.
					err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					if err != nil {
						log.Println("write close:", err)
						return
					}
					select {
					case <-done:
					case <-time.After(time.Second):
					}
					return
				}
			}
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
