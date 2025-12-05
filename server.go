package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

type Message struct {
	From      string    `json:"from"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // msg | system
}

type Client struct {
	id   string
	conn net.Conn
	send chan Message
}

type Server struct {
	mu      sync.Mutex
	clients map[string]*Client
	nextID  int
}

func NewServer() *Server {
	return &Server{
		clients: make(map[string]*Client),
		nextID:  1,
	}
}

func (s *Server) broadcast(msg Message, exceptID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, c := range s.clients {
		if id == exceptID {
			continue // no self-echo
		}
		select {
		case c.send <- msg:
		default:
			log.Println("Dropped message to client", id)
		}
	}
}

func (s *Server) addClient(conn net.Conn) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := strconv.Itoa(s.nextID)
	s.nextID++

	client := &Client{
		id:   id,
		conn: conn,
		send: make(chan Message, 32),
	}

	s.clients[id] = client
	return client
}

func (s *Server) removeClient(c *Client) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.clients, c.id)
	close(c.send)
}

func (s *Server) handleClient(c *Client) {
	// Goroutine لإرسال الرسائل للعميل
	go func() {
		writer := bufio.NewWriter(c.conn)
		enc := json.NewEncoder(writer)

		for msg := range c.send {
			enc.Encode(msg)
			writer.Flush()
		}
	}()

	// إشعار باقي المستخدمين بدخول المستخدم
	joinMsg := Message{
		From:      "system",
		Text:      fmt.Sprintf("User [%s] joined", c.id),
		Timestamp: time.Now(),
		Type:      "system",
	}
	s.broadcast(joinMsg, "")

	// استقبال الرسائل من العميل
	dec := json.NewDecoder(c.conn)

	for {
		var msg Message
		if err := dec.Decode(&msg); err != nil {
			break
		}

		msg.From = c.id
		msg.Timestamp = time.Now()
		msg.Type = "msg"

		// بث الرسالة لكل المستخدمين ماعدا المرسل
		s.broadcast(msg, c.id)
	}

	// عند الخروج
	s.removeClient(c)

	leaveMsg := Message{
		From:      "system",
		Text:      fmt.Sprintf("User [%s] left", c.id),
		Timestamp: time.Now(),
		Type:      "system",
	}
	s.broadcast(leaveMsg, "")

	c.conn.Close()
}

func main() {
	port := os.Getenv("CHAT_PORT")
	if port == "" {
		port = "1234"
	}

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Real-time chat server listening on port", port)

	server := NewServer()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		client := server.addClient(conn)
		go server.handleClient(client)
	}
}
