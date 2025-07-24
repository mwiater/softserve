// serve.go
package softserve

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"net"
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

// StartServerInternalCerts starts the server, generating self-signed certificates
// in memory if SSL is enabled, instead of loading them from files.
func StartServerInternalCerts() error {
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
		// Generate internal self-signed certs
		cert, err := GenerateInternalSelfSignedCert()
		if err != nil {
			return fmt.Errorf("failed to generate internal self-signed cert: %w", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		server.Addr = fmt.Sprintf(":%d", AppConfig.HTTPSPort)
		server.TLSConfig = tlsConfig // Assign the TLS config to the server

		go func() {
			fmt.Printf("🔒 Serving HTTPS on https://0.0.0.0:%d (in-memory certs)\n", AppConfig.HTTPSPort)
			// When TLSConfig is set on the server, ListenAndServeTLS can be called with empty paths
			if err := server.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
				log.Fatalf("Server error with internal certs: %v", err)
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

// GenerateSelfSignedCert generates a self-signed TLS certificate and private key in memory.
// It returns a tls.Certificate suitable for direct use in tls.Config.
func GenerateInternalSelfSignedCert() (tls.Certificate, error) {
	// Generate a new RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate RSA key: %w", err)
	}

	// Create a certificate template
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"SoftServe Internal"},
			CommonName:   "localhost",
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // Valid for 1 year

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false, // Not a CA, just a server cert
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")}, // Include IPv6 for localhost
	}

	// Create the self-signed certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Load the certificate and key into a tls.Certificate struct
	// Note: tls.X509KeyPair expects PEM-encoded data. For in-memory, we directly construct.
	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  privateKey,
	}

	return cert, nil
}
