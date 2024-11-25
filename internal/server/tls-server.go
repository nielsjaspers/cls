package server

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/nielsjaspers/cls/secrets"
)

func SetupTLSServer() {
	log.SetFlags(log.Lshortfile)

	certificate, err := tls.LoadX509KeyPair(secrets.ServerCrtPath, secrets.ServerKeyPath)
	if err != nil {
		log.Printf("Failed to load Keypair: %v", err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{certificate}}
	config.MinVersion = tls.VersionTLS12
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Printf("Failed to open tls port: %v", err)
		return
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Connection failed: %v", err)
			continue // try again
		}

		tlsConn, ok := conn.(*tls.Conn)
		if ok {
			err := tlsConn.Handshake()
			if err != nil {
				log.Printf("TLS handshake failed: %v", err)
				continue
			}

			state := tlsConn.ConnectionState()
			log.Printf("TLS handshake complete: %v, Version: %v, CipherSuite: %v", state.HandshakeComplete, state.Version, state.CipherSuite)
		} else {
			log.Println("Received non-TLS connection")
		}
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Debugging
	log.Println("Client connected")

	// Specify the file path 
    // Currently hardcoded as this is a proof of concept
	filePath := "~/received/received_file.jpeg" 

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting home directory: %v", err)
		return
	}
	filePath = strings.Replace(filePath, "~", homeDir, 1)

	// Open file for writing (create or truncate)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	// Read file content from connection and write it to the file
	r := bufio.NewReader(conn)
	buf := make([]byte, 4096) // 4KB chunks
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("File transfer complete")
				break
			}
			log.Printf("Error reading: %v", err)
			return
		}

		if _, err := file.Write(buf[:n]); err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}
	}

	log.Printf("File successfully saved to %s", filePath)

	// Respond to the client
	_, err = conn.Write([]byte("File received successfully!\n"))
	if err != nil {
		log.Printf("Error writing to client: %v", err)
	}
}

