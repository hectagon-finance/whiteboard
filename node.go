package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	// "github.com/onspeedhp/project/crypto"
)

// Node contains server and client
var upgrader = websocket.Upgrader{} // use default options

var (
	MessageString = 1
	MessagePeers  = 2
	MessageInput = 3
)

type MemPool struct {
	Inputs []*Input `json:"transactions"`
}

type Input struct { 
	Id 		  string `json:"id"`
	PublicKey string `json:"publicKey"`
	Signature string `json:"signature"`
	Data []string 	 `json:"data"`
}

type Message struct {
	From    string   `json:"from"`
	To      string   `json:"to"`
	Content []interface{} `json:"content"`
	Type    int		`json:"type"`
}

// type Node interface {
// 	Serve(w http.ResponseWriter, r *http.Request)
// 	// Start the node's server
// 	StartServer(port string)
// 	// Start the node's client
// 	StartClient()
// 	// Get the node's address
// 	Address() string
// 	// Get the node's port
// 	Port() string
// 	// Get the node's id
// 	Id() string
// 	// Connect to another node
// 	ConnectToNode(node *Node)
// }

type node struct {
	address string
	port    string
	id      string
	peers   []interface{}
	memPool MemPool
	
}

func (node *node) Address() string {
	return node.address
}

func (node *node) Port() string {
	return node.port
}

func (node *node) Id() string {
	return node.id
}

func (n *node) Serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			// log.Println("read:", err)
			break
		}
		log.Printf("Server: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Println("unmarshal:", err)
			break
		}

		if msg.Type == 1 {
			log.Println("MessageString")
		} else if msg.Type == 2 {
			if CompareSlices(msg.Content, n.peers) {
				// log.Println("Peers already exist")
			} else {
				n.AddToPeers(msg.Content)
			}
		} else if msg.Type == 3 {
			time.Sleep(1 * time.Second)
			input := ConvertInterfacetoInput(msg.Content[0])
			n.AddtoMemPool(&input)
			n.Broadcast(msg)
			// exist := n.AddtoMemPool(&input)
			// // n.Broadcast(msg)
			// if exist {
			// 	// log.Println("Input already exist")
			// } else {
			// 	n.Broadcast(msg)
			// }
			log.Println(len(n.memPool.Inputs))
		}
	}
}

func (node *node) Start() {
	log.Println("Start server of node at " + node.port)
	http.HandleFunc("/ws", node.Serve)
	log.Fatal(http.ListenAndServe(":"+node.port, nil))
}

func (n *node) StartClient() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: ":" + n.port, Path: "/ws"}
	log.Println("Start client of node at " + n.port)

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer c.Close()

	done := make(chan struct{})
	
	go func() {
		defer close(done)
		for {
			// _, message, err := c.ReadMessage()
			// if err != nil {
			// 	return
			// }
			// log.Printf("Client: %s", message)
		}
	}()

	// go n.ConnectToNode()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		// case t := <-ticker.C:
		// 	err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		// 	if err != nil {
		// 		return
		// 	}
		// }
		case <-ticker.C:
			go n.Broadcast()	
		}
	}
}

// User client to connect to another node
func (n *node) Broadcast(args ...interface{}) {
	if len(n.peers) != 0 {
		for _, port := range n.peers {
			if n.port != port {
				peer := &node{address: "localhost", port: port.(string), id: "0"}
				// connect to the peer's server and send a message
				u := url.URL{Scheme: "ws", Host: peer.Address() + ":" + peer.Port(), Path: "/ws"}

				c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
				if err != nil {
					// log.Println("dial:", err)
					return
				}
				defer c.Close()
				
				if len(args) == 1 {
					message := args[0]
					// convert message to Message type
					msg := message.(Message)

					msg.From = n.port
					msg.To = peer.port
					time.Sleep(1 * time.Second)
					c.WriteJSON(msg)
				} else if len(args) > 1 {
					return
				} else {
					message := Message{From: n.port, To: peer.port, Content: n.peers, Type: MessagePeers}
					c.WriteJSON(message)
				}

			}
		}
	}
}

func ConvertStringtoInterface(str []string) []interface{} {
	var result []interface{}
	for _, v := range str {
		result = append(result, v)
	}
	return result
}

func ConvertInterfacetoInput(inter interface{}) Input {
	// Assert msg is a map[string]interface{}
	msgMap, ok := inter.(map[string]interface{})
	if !ok {
		// Handle error: msg is not a map[string]interface{}
	}

	// Extract values from msgMap and assign to new instance of Input
	input := Input{
		Id : 	   msgMap["id"].(string),
		PublicKey: msgMap["publicKey"].(string),
		Signature: msgMap["signature"].(string),
		Data:      []string{msgMap["data"].(string)},
	}

	return input
}

func (n *node) AddToPeers(peers []interface{}) []interface{}{
    // check for new peers
    for _, peer := range peers {
		alreadyExists := false
		for _, p := range n.peers {
			if p == peer {
				alreadyExists = true
				break
			}
		}
		if !alreadyExists {
			n.peers = append(n.peers, peer)
		}
    }
	return n.peers
}

func (n *node) AddtoMemPool(in *Input) bool{
	alreadExists := false
	for _, input := range n.memPool.Inputs {

		if input.Id == in.Id {
			alreadExists = true
			break
		}
	}

	if !alreadExists {
		n.memPool.Inputs = append(n.memPool.Inputs, in)
	}

	return alreadExists
}

func (n *node) GeneratePeers(){
	n.peers = append(n.peers, "8080")
	if n.port != "8080" {
		n.peers = append(n.peers, n.port)
	}
}

func CompareSlices(a, b []interface{}) bool {
    if len(a) != len(b) {
		// log.Println("Slices are not the same length")
        return false
    }

    for i := 0; i < len(a); i++ {
        if !compareValues(a[i], b[i]) {
            return false
        }
    }

    return true
}

func compareValues(a, b interface{}) bool {
    // Handle nil values
    if a == nil && b == nil {
        return true
    } else if a == nil || b == nil {
        return false
    }
    // Type assertion
    switch a.(type) {
    case int:
        return a.(int) == b.(int)
    case string:
        return a.(string) == b.(string)
    case bool:
        return a.(bool) == b.(bool)
    // Add other types as necessary...
    default:
        // Not comparable types
        return false
    }
}


func main(){
	port := os.Args[1]

	// create a new node
	node := node{
		address: "localhost",
		port: port,
		id: "1",
	}
	node.GeneratePeers()
	// start the node's server
	go node.Start()

	// start the node's client
	go node.StartClient()

	for {
	}
}