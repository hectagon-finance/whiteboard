# Welcome to whiteboard chain

## What we build ?

- We build an appchain with simple logic like a todo app to help everyone understand how a blockchain works

## How to run ?

### Step 1: Run validator

- Let's take a look in the Makefile.
- You can run validator by using make command or using direct command.

Example:

```
Note: First you need to run the validator at port 8080
Try: "make run1" or "go run main.go 8080 genesis"
Open another terminal and continue:
     "make run2" or "go run main.go 9000 8080"
```

You can continue to run more validators as you like.

## Step 2: Run client

- Go to the folder client: cd cmd/client

* Note: Don't care about the client folder in the root folder we are developing an interface to make it easier for you to see. You only need to pay attention to the client folder in the cmd folder

- In cmd/client folder we have keypair file for you testing. Now you can use client to interact with validator. We are defaulting the client to send the transaction to the validator whose port is 9000. Ok let's try:

```
Try "make client1" or
"go run main.go send "Hello" -k 8df4135ecefc9a4d054e2c596cd3f56432e683431b27216fea917b01c8ef1fee"
```

- after -k is the user's private key. This private key and the public key in file public_key.txt need to be a key pair.
