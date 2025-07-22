// serve.go
package softserve

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

func StartServer() error {
	hub := newReloadHub()

	mux := http.NewServeMux()

	// WebSocket route
	mux.Handle("/__ws", websocket.Handler(func(ws *websocket.Conn) {
		hub.mu.Lock()
		hub.clients[ws] = true
		hub.mu.Unlock()

		// Block until disconnect
		buf := make([]byte, 1)
		ws.Read(buf)

		hub.mu.Lock()
		delete(hub.clients, ws)
		hub.mu.Unlock()
	}))

	// File server route
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveFileWithInjection(w, r)
	})

	go watchForChanges(AppConfig.WebRoot, hub)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", AppConfig.HTTPPort),
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		fmt.Printf("🌐 Serving %s on http://0.0.0.0:%d\n", AppConfig.WebRoot, AppConfig.HTTPPort)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("🛑 Shutting down...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return server.Shutdown(shutdownCtx)
}

func serveFileWithInjection(w http.ResponseWriter, r *http.Request) {
	relPath := strings.TrimPrefix(r.URL.Path, "/")
	if relPath == "" {
		relPath = "index.html"
	}
	fullPath := filepath.Join(AppConfig.WebRoot, relPath)

	if strings.HasSuffix(fullPath, ".html") {
		data, err := os.ReadFile(fullPath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		modified := injectReloadScript(string(data))
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(modified))
		return
	}

	http.ServeFile(w, r, fullPath)
}

func injectReloadScript(html string) string {
	script := `<script>
		const ws = new WebSocket("ws://" + location.host + "/__ws");
		ws.onmessage = () => location.reload();
	</script>`

	if strings.Contains(html, "</body>") {
		return strings.Replace(html, "</body>", script+"</body>", 1)
	}
	return html + script
}
