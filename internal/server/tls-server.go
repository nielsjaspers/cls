package server

import (
	"bufio"
	"crypto/tls"
	"io"
	"log"
	"net"

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
				return
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

	// debugging
	log.Println("Client connected")

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

		_, err = conn.Write([]byte("Message recieved, hello client!\n"))
		if err != nil {
			log.Printf("Error writing: %v", err)
			return
		}
	}
}
