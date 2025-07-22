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
)

func StartServer() error {
	hub := newReloadHub()

	mux := http.NewServeMux()

	// WebSocket endpoint
	mux.Handle("/__ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hub.handleWebSocket(w, r)
	}))

	// Static + API + injected HTML handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if HandleAPIRequest(w, r) {
			return
		}
		serveFileWithInjection(w, r)
	})

	go watchForChanges(AppConfig.WebRoot, hub)

	server := &http.Server{
		Addr:    "", // Set below
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	if AppConfig.SSL {
		// Default fallback
		if AppConfig.CertsPath == "" {
			AppConfig.CertsPath = "certs"
		}

		certPath := filepath.Join(AppConfig.CertsPath, "cert.pem")
		keyPath := filepath.Join(AppConfig.CertsPath, "key.pem")

		if AppConfig.GenerateCerts {
			if err := GenerateSelfSignedCert(AppConfig.CertsPath); err != nil {
				return fmt.Errorf("cert generation failed: %w", err)
			}
		}

		fmt.Println("SSL: Loading Cert files:")
		fmt.Println("  >>>", certPath)
		fmt.Println("  >>>", keyPath)

		server.Addr = fmt.Sprintf(":%d", AppConfig.HTTPSPort)
		go func() {
			fmt.Printf("🔒 Serving HTTPS on https://0.0.0.0:%d\n", AppConfig.HTTPSPort)
			if err := server.ListenAndServeTLS(certPath, keyPath); err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}

		}()
	} else {
		server.Addr = fmt.Sprintf(":%d", AppConfig.HTTPPort)
		go func() {
			fmt.Printf("🌐 Serving HTTP on http://0.0.0.0:%d\n", AppConfig.HTTPPort)
			if err := server.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatalf("Server error: %v", err)
			}
		}()
	}

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
	wsScheme := "ws"
	if AppConfig.SSL {
		wsScheme = "wss"
	}

	script := fmt.Sprintf(`<script>
		const ws = new WebSocket("%s://" + location.host + "/__ws");
		ws.onmessage = () => location.reload();
	</script>`, wsScheme)

	if strings.Contains(html, "</body>") {
		return strings.Replace(html, "</body>", script+"</body>", 1)
	}
	return html + script
}
