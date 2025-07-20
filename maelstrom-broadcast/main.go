package main

import (
	"encoding/json"
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	receivedMessages := make([]float64, 0)	// kept separate per each dest node
	var topology map[string] any

	n := maelstrom.NewNode()

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		// Unmarshal message body
		var body map[string] any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		// Record message in receivedMessages
		if msgInt, ok := body["message"].(float64); ok {
			receivedMessages = append(receivedMessages, msgInt)
		} else {
			return fmt.Errorf("unable to add message integer: message is type %T", body["message"])
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
			"messages": receivedMessages,
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