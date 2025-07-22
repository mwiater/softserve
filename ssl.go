// ssl.go
package softserve

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"
)

func GenerateSelfSignedCert(path string) error {
	certPath := filepath.Join(path, "cert.pem")
	keyPath := filepath.Join(path, "key.pem")

	fmt.Printf("Checking for existing cert path: '%s'\n", path)
	if err := ensureAbsoluteAndExists(path); err != nil {
		fmt.Printf("  Error: %v\n\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("  Success: Path is an absolute, existing directory.\n")
	}

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %w", err)
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		return fmt.Errorf("failed to generate serial number: %w", err)
	}

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			Organization: []string{"Softserve Dev Cert"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	der, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return fmt.Errorf("failed to create cert: %w", err)
	}

	// Write cert
	certOut, err := os.Create(certPath)
	if err != nil {
		return err
	}
	defer certOut.Close()
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: der})

	// Write key
	keyOut, err := os.Create(keyPath)
	if err != nil {
		return err
	}
	defer keyOut.Close()
	pem.Encode(keyOut, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	fmt.Printf("🔐 Generated self-signed cert at %s/\n", path)
	return nil
}

// ensureAbsoluteAndExists checks if the given path is an absolute path
// and if it already exists as a directory. It returns an error if
// either condition is not met.
func ensureAbsoluteAndExists(path string) error {
	// 1. Check if the path is absolute
	if !filepath.IsAbs(path) {
		return fmt.Errorf("path is not absolute: %s", path)
	}

	// 2. Check if the path exists and is a directory
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Path does not exist, which is the desired error condition
			return fmt.Errorf("path does not exist: %s", path)
		}
		// Some other error occurred while trying to stat the path (e.g., permissions)
		return fmt.Errorf("failed to stat path %s: %w", path, err)
	}

	// 3. Check if the existing path is a directory
	if !fileInfo.IsDir() {
		return fmt.Errorf("path exists but is not a directory: %s", path)
	}

	// If all checks pass, the path is an absolute, existing directory
	return nil
}
