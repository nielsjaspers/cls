package client

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/nielsjaspers/cls/secrets"
)

func SetupTLSClient() {
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

	_, err = conn.Write([]byte("Hello Server!\n"))
	if err != nil {
		log.Printf("Error writing: %v", err)
		return
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

	// TODO:    Wait before disconnecting from server
	//          Have option to send message to server -- (and to get a response)

	// TODO2:   Function to send file to server
	//          Function to import file from server

}

// fileUpload changes file f to []byte for sending over tls
// returns a byte array with file content, and an error if there is one 
func fileUpload(f *os.File) ([]byte, error) {
    fmt.Printf("%v", f)
    return []byte(""), nil
}

// getRemoteFilePaths requests all remote paths of files currently on the server
// returns a string array with paths to the files, and an error if there is one
func getRemoteFilePaths() ([]string, error) { // Use TLS connection (conn) maybe as parameter ?
    return []string{"",""}, nil
}

// getRemoteFile requests a single file with remote path p
// returns a file and an error if there is one
func getRemoteFile(p string) (os.File, error) {
    fmt.Printf("%v", p)
    return *os.NewFile(1, ""), nil
}
