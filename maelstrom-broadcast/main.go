package main

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"slices"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	receivedMessages := make(map[float64] struct{}, 0)	// kept separate per each dest node
	var topology map[string] any

	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal message body
		var body map[string] any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		msgInt, ok := body["message"].(float64);
		if !ok {
			return fmt.Errorf("unable to parse message as integer: message is type %T", body["message"])
		}
		
		// Check if the message is new
		if _, exists := receivedMessages[msgInt]; !exists {
			// Record message in receivedMessages
			receivedMessages[msgInt] = struct{}{}
			
			// Broadcast to neighbors not in the sender's network
			broadcast_list := make(map[string] struct{}, 0)
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
					// Unmarshal message body
					var body map[string] any
					if err := json.Unmarshal(msg.Body, &body); err != nil {
						return err
					}

					// Ensure reply from recipients
					if msgtype, ok := body["type"].(string); (ok || msgtype != "broadcast_ok") {
						return fmt.Errorf("did not get ok from %v", node)
					}

					return nil
				})
			}
		}

		// Create response body
		response_body := map[string] any {
			"type": "broadcast_ok",
		}

		return n.Reply(msg, response_body)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		body := map[string] any {
			"type": "read_ok",
			"messages": slices.Collect(maps.Keys(receivedMessages)),
		}

		return n.Reply(msg, body)
	})

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string] any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Set topology if received for the first time
		if topology == nil {
			topology = body["topology"].(map[string] any)
		}

		request_body := map[string] any {
			"type": "topology_ok",
		}

		return n.Reply(msg, request_body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}