package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

type Message struct {
	From      string    `json:"from"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "msg" | "system"
}

func readServer(conn net.Conn, stop chan struct{}) {
	dec := json.NewDecoder(conn)
	for {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			log.Printf("connection closed by server or decode error: %v", err)
			close(stop)
			return
		}
		// pretty print
		if msg.Type == "system" {
			fmt.Printf("\n[system] %s\n", msg.Text)
		} else {
			fmt.Printf("\n[%s] %s\n", msg.From, msg.Text)
		}
		fmt.Print("> ")
	}
}

func main() {
	addr := os.Getenv("CHAT_ADDR")
	if addr == "" {
		addr = "localhost:1234"
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("failed to connect to server at %s: %v", addr, err)
	}
	defer conn.Close()

	fmt.Println("Connected to chat server at", addr)
	stop := make(chan struct{})
	go readServer(conn, stop)

	// writer uses encoder
	w := bufio.NewWriter(conn)
	enc := json.NewEncoder(w)

	console := bufio.NewReader(os.Stdin)
	fmt.Println("Type messages (type 'exit' to quit):")
	for {
		fmt.Print("> ")
		line, _ := console.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.EqualFold(line, "exit") {
			fmt.Println("Bye!")
			return
		}
		msg := Message{
			Text:      line,
			Timestamp: time.Now(),
			Type:      "msg",
		}
		if err := enc.Encode(msg); err != nil {
			log.Printf("encode error: %v", err)
			return
		}
		if err := w.Flush(); err != nil {
			log.Printf("flush error: %v", err)
			return
		}

		// check if server closed
		select {
		case <-stop:
			return
		default:
		}
	}
}
