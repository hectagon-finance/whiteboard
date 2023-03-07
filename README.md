# whiteboard
Build a blockchain from scratch to understand what compromise a chain-designer have to make

# How a blockchain works?
## There are 5 routines:
- FindPeer
- ClientHandler
- Gossip
- Consensus
- Logic

## There are 6 channels:
- C1 chan Peer. FindPeer -> Gossip
- C2 chan Peer. FindPeer -> Consensus
- C3 chan Mempool. Gossip -> Consesus
- C4 chan BlockHash. Consensus -> Gossip
- C5 chan []Input. ClientHandler -> Gossip
- C6 chan Block. Consensus -> Logic

## Description
