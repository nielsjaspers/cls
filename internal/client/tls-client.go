package client

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"io"
	"log"
	"os"

	"github.com/nielsjaspers/cls/secrets"
)

func SetupTLSClient(f []byte) {
	log.SetFlags(log.Lshortfile)

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

	conn, err := tls.Dial("tcp", secrets.ServerURL, conf)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to server")

	// _, err = conn.Write([]byte("Hello Server!\n"))
	// if err != nil {
		// log.Printf("Error writing: %v", err)
		// return
	// }

    
    if len(f) > 0 {
        _, err = conn.Write(f)
        if err != nil {
            log.Printf("Error writing: %v", err)
            return
        }
    }

	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Client disconnected")
			} else {
				log.Printf("Error reading: %v", err)
			}
			return
		}

		log.Printf("Received: %s", msg)
	}

}



