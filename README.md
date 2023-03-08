# whiteboard

Build a blockchain from scratch to understand what compromise a chain-designer have to make

# How a blockchain works?

## There are 5 routines:

- PeerDiscovery
- ClientWriteHandler
- MessageGossiper
- Consensus
- ChainLogic

## There are 6 channels:

- C1 chan Peer. PeerDiscovery -> MessageGossiper
- C2 chan Peer. PeerDiscovery -> Consensus
- C3 chan Mempool. MessageGossiper -> Consesus
- C4 chan BlockHash. Consensus -> MessageGossiper
- C5 chan []Input. ClientWriteHandler -> MessageGossiper
- C6 chan BlockData. Consensus -> ChainLogic

## Data type

- Peer
- Input
- Block
- BlockHash
- Mempool

## Description

1. How do nodes in a peer to peer network work with each other?

- Nodes can find each other through the list of peers at PeerDiscovery.
- When any node joins or leaves the network, the list of peers will be updated at PeerDiscovery. PeerDiscovery will pass a constantly updated list of nodes to MessageGossiper and Consensus through chan C1 and chan C2.

2. What happens when the client interact with the blockchain network?

- When the client write request (An input include: {pubk, raw, hashed}) to the blockchain network
  => ClientWriteHandler will validate input and send to MessageGossiper through chan C5.

3. Consensus Mechanism

- After receiving the input, node will add input to the mempool and sort the inputs in order of higher gas fee will be above. Then MessageGossiper will broadcast the input to the entire network based on the list of nodes that PeerDiscovery provided.
- When the first node collects enough input for a block, MessageGossiper will pack it and send it to Consensus through chan C3 and broadcast to the network based on the list of nodes that PeerDiscovery provided.
- Consensus validate inputs, generates a hash of the block and broadcast the blockhash to the network based on the list of nodes that PeerDiscovery provided. Consensus also received results from other nodes. Blockhash with the largest weight will be considered valid, the block with valid hash will be generated and add to chain. Then the block will be sent to Chain logic through chan C6.
- The valid blockhash will be sent by Consensus to MessageGossiper through chan C4. Nodes with different blockhash will not be able to continue communicating with each other.
  Node want to continue communicating will have to synchronize with the latest chain.

4. Logic of Appchain

- The input will be valid if it is in a valid block. From the valid input we can write functions to handle the appchain logic and send results to client.
