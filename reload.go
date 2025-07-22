// reload.go
package softserve

import (
	"log"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type reloadHub struct {
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
}

func newReloadHub() *reloadHub {
	return &reloadHub{
		clients: make(map[*websocket.Conn]bool),
	}
}

func (h *reloadHub) run() {
	http.Handle("/__ws", websocket.Handler(func(ws *websocket.Conn) {
		h.mu.Lock()
		h.clients[ws] = true
		h.mu.Unlock()

		// Block until disconnected
		buf := make([]byte, 1)
		ws.Read(buf)

		h.mu.Lock()
		delete(h.clients, ws)
		h.mu.Unlock()
	}))
}

func (h *reloadHub) broadcastReload() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if _, err := conn.Write([]byte("reload")); err != nil {
			log.Printf("removing stale client: %v", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}
