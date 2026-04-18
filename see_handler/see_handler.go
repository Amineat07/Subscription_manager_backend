package ssehandler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type SSEMessage struct {
	Event string
	Data  []byte
}

var (
	clients = make(map[chan SSEMessage]bool)
	mu      sync.RWMutex
)

func NewsFeedSSE(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	ch := make(chan SSEMessage, 10)

	mu.Lock()
	clients[ch] = true
	mu.Unlock()

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer func() {
			mu.Lock()
			delete(clients, ch)
			close(ch)
			mu.Unlock()
		}()

		fmt.Fprintf(w, ": connected\n\n")
		w.Flush()

		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case msg, ok := <-ch:
				if !ok {
					return
				}
				fmt.Fprintf(w, "event: %s\n", msg.Event)
				fmt.Fprintf(w, "data: %s\n\n", msg.Data)
				if err := w.Flush(); err != nil {
					return
				}
			case <-ticker.C:
				fmt.Fprintf(w, ": heartbeat\n\n")
				if err := w.Flush(); err != nil {
					return
				}
			}
		}
	})

	return nil
}


func BroadcastNewsFeed(data interface{}, eventType string) {
	payload, err := json.Marshal(data)
	if err != nil {
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	for ch := range clients {
		select {
		case ch <- SSEMessage{Event: eventType, Data: payload}:
		default:
		}
	}
}