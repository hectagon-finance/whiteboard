package validator

import (
	"encoding/hex"
	"log"

	"github.com/hectagon-finance/whiteboard/types"
)

func BroadcastBlockHash() {
	for {
		select {

		case msg := <-Chan_1:

			log.Print("** routine2: Client Handler** bool:", (DraftBlock.Hash == [32]byte{}) && (msg.memPool.Size() >= msg.memPool.CutOff))

			if (DraftBlock.Hash == [32]byte{}) && (msg.memPool.Size() >= msg.memPool.CutOff) {
				var k int

				if msg.memPool.Size() >= msg.memPool.CutOff {
					k = msg.memPool.CutOff
				} else {
					k = msg.memPool.Size()
				}

				DraftBlock = types.NewBlock(msg.heigt, [32]byte{}, msg.memPool.Transactions[:k])

				// broadcast block hash
				blockHashSlice := DraftBlock.Hash[:]
				blockHashStr := hex.EncodeToString(blockHashSlice)

				message := map[string]interface{}{
					"type":      "blockHash",
					"from":      Port,
					"blockHash": blockHashStr,
				}
				ShouldReceiveTxFromPeer = false
				Chan_2 <- k

				ConnectAndSendMessage(message)
			}
		}
	}
}
