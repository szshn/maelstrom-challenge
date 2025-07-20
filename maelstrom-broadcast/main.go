package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	var receivedMessages sync.Map	// kept separate per each dest node
	var topology map[string]any

	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		return broadcast_handler(msg, &receivedMessages, topology, n)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		return read_handler(msg, &receivedMessages, n)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		return topology_handler(msg, &topology, n)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}

func broadcast_handler(
	msg maelstrom.Message,
	receivedMessages *sync.Map,
	topology map[string]any,
	n *maelstrom.Node,
) error {
	// Unmarshal message body
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	msgInt, ok := body["message"].(float64);
	if !ok {
		return fmt.Errorf("unable to parse message: message is type %T", body["message"])
	}
	
	// Check if the message is new
	if _, exists := receivedMessages.LoadOrStore(msgInt, struct{}{}); !exists {
		// Broadcast to neighbors not in the sender's network
		broadcast_list := make(map[string]struct{}, len(topology[msg.Dest].([]any)))	// set length as length of toplogy[msg.Dest]?
		if neighbors, ok := topology[msg.Dest].([]any); ok {
			for _, node := range neighbors {
				if node_name, ok := node.(string); ok {
					broadcast_list[node_name] = struct{}{}
				}
			}
		}
		if neighbors, ok := topology[msg.Src].([]any); ok {
			for _, node := range neighbors {
				if node_name, ok := node.(string); ok {
					delete(broadcast_list, node_name)
				}
			}
		}
		for node := range broadcast_list {
			n.RPC(node, body, func(msg maelstrom.Message) error {
				return nil
			})
		}
	}

	// Create response body
	response_body := map[string]any {
		"type": "broadcast_ok",
	}

	return n.Reply(msg, response_body)
}

func read_handler(
	msg maelstrom.Message,
	receivedMessages *sync.Map,
	n *maelstrom.Node,
) error {
	messages := make([]float64, 0)
    receivedMessages.Range(func(key, value any) bool {
		messages = append(messages, key.(float64))
        return true
    })

	body := map[string]any {
		"type": "read_ok",
		"messages": messages,
	}

	return n.Reply(msg, body)
}

func topology_handler(
	msg maelstrom.Message,
	topology *map[string]any,
	n *maelstrom.Node,
) error {
	var body map[string]any
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		return err
	}

	// Set topology if received for the first time
	if *topology == nil {
		if topo, ok := body["topology"].(map[string]any); ok {
			*topology = topo
		} else {
			return fmt.Errorf("invalid topology format")
		}
	}

	request_body := map[string]any{
		"type": "topology_ok",
	}

	return n.Reply(msg, request_body)
}