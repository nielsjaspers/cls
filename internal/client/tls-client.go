package client

import (
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/nielsjaspers/cls/secrets"
)

func SetupTLSClient() {
    log.SetFlags(log.Lshortfile)

    // Lees het CA-certificaat in
    cert, err := os.ReadFile(secrets.CertAuthPath)
    if err != nil {
        log.Fatalf("Failed to read certificate file: %v", err)
    }

    caCertPool := x509.NewCertPool()
    if !caCertPool.AppendCertsFromPEM(cert) {
        log.Fatalf("Failed to append cert to pool")
    }

    conf := &tls.Config{
        RootCAs: caCertPool,
    }

    // Maak een TLS-verbinding met de server
    conn, err := tls.Dial("tcp", "localhost:443", conf)
    if err != nil {
        log.Fatalf("Failed to connect to server: %v", err)
    }
    defer conn.Close()

    log.Println("Connected to server")
}
