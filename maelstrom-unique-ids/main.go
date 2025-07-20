package main

import (
	"fmt"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	unique_id := 0

	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		body := map[string] any {
			"type": "generate_ok",
			"id": fmt.Sprintf("%v%v%d", msg.Src, msg.Dest, unique_id),
		}
		
		unique_id += 1

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}