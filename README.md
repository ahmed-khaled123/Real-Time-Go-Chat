# Real-Time Go Chat (Assignment 04)

A real-time broadcasting chat system built with Go using **goroutines**, **channels**, and **mutex** for concurrency control.  
This project is an upgrade from the RPC-based chat to a fully concurrent real-time system.

---

## ✅ Features
- ✅ Real-time message broadcasting
- ✅ User join notification: `User [ID] joined`
- ✅ No self-echo (sender does not receive their own message)
- ✅ Concurrent send/receive using goroutines and channels
- ✅ Shared client list protected using mutex
- ✅ Multiple clients supported simultaneously
- ✅ TCP-based communication (no RPC)

---

## 🧠 Technologies Used
- Go (Golang)
- Goroutines
- Channels
- sync.Mutex
- TCP Networking
- JSON Encoding/Decoding
- (Optional) Docker

---

## 📂 Project Structure
realtime-go-chat/
├── server.go # Real-time chat server
├── client.go # Terminal chat client
├── Dockerfile # (Optional) Server Dockerfile
└── README.md # Project documentation

## ▶️ How to Run Locally

### 1️⃣ Run the server:
```bash
go run server.go
You should see: Real-time chat server listening on port 1234


2️⃣ Run clients (in separate terminals): go run client.go

Each new client will trigger: [system] User [ID] joined


