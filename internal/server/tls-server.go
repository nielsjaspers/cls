package server

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"os"

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

	r := bufio.NewReader(conn)

	//Listen for filename (max 255 bytes)
	filenameBuf := make([]byte, 255)
	n, err := r.Read(filenameBuf)
	if err != nil {
		log.Printf("Error reading filename: %v", err)
		return
	}
	filename := string(bytes.Trim(filenameBuf[:n], "\x00\n"))
	log.Printf("Received filename: %s", filename)

	// Respond with "NEXT_ITEM"
	_, err = conn.Write([]byte("NEXT_ITEM\n"))
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}

	// Listen for file extension (max 15 bytes)
	extensionBuf := make([]byte, 15)
	n, err = r.Read(extensionBuf)
	if err != nil {
		log.Printf("Error reading file extension: %v", err)
		return
	}
	extension := string(bytes.Trim(extensionBuf[:n], "\x00\n"))
	log.Printf("Received file extension: %s", extension)

	// Respond with "NEXT_ITEM"
	_, err = conn.Write([]byte("NEXT_ITEM\n"))
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}

	filePath := fmt.Sprintf("%s", filename)

	// Open the file for writing
	homeDir, err := os.UserHomeDir()
	file, err := os.Create(homeDir + "/" + filePath)
	if err != nil {
		log.Printf("Error creating file: %v", err)
		return
	}
	defer file.Close()

	// Listen for file content (no max size)
	buf := make([]byte, 131072) // 128 kB chunks
	for {
		n, err := r.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Println("File transfer complete")
				break
			}
			log.Printf("Error reading file content: %v", err)
			return
		}

		// Check if the received data contains the EOF marker
		if bytes.Contains(buf[:n], []byte("EXIT_EOF_EXIT_EOF\n")) {
			log.Println("EOF marker received, file transfer complete")
			break
		}

		// Write the received content to the file
		if _, err := file.Write(buf[:n]); err != nil {
			log.Printf("Error writing to file: %v", err)
			return
		}
	}

	log.Printf("File successfully saved to %s", filePath)

	// Respond to the client after the file is fully received
	_, err = conn.Write([]byte("File received successfully!\n"))
	if err != nil {
		log.Printf("Error sending final response: %v", err)
	}
}
