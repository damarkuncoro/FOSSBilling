package notifications

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Relax for dev, tighten in production if needed
}

type WSHub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

func NewWSHub() *WSHub {
	return &WSHub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *WSHub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WSHub] Upgrade failed: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	log.Printf("[WSHub] Client connected. Total clients: %d", len(h.clients))

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		conn.Close()
		log.Printf("[WSHub] Client disconnected. Total clients: %d", len(h.clients))
	}()

	// Keep connection alive, listen for nothing for now
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("[WSHub] ReadMessage error: %v", err)
			break
		}
	}
}

func (h *WSHub) Broadcast(ctx context.Context, msg any) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if err := client.WriteMessage(websocket.TextMessage, payload); err != nil {
			log.Printf("WS Broadcast error: %v", err)
		}
	}

	return nil
}
