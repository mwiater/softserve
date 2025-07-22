// reload.go
package softserve

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type reloadHub struct {
	clients  map[*websocket.Conn]bool
	mu       sync.Mutex
	upgrader websocket.Upgrader
}

func newReloadHub() *reloadHub {
	return &reloadHub{
		clients: make(map[*websocket.Conn]bool),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // allow all origins
			},
		},
	}
}

func (h *reloadHub) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WS upgrade error: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, conn)
		h.mu.Unlock()
		conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			break
		}
	}
}

func (h *reloadHub) broadcastReload() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, []byte("reload")); err != nil {
			log.Printf("removing stale WS client: %v", err)
			conn.Close()
			delete(h.clients, conn)
		}
	}
}
