package server

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/nielsjaspers/cls/internal/arguments"
	filehandler "github.com/nielsjaspers/cls/pkg"
	"github.com/nielsjaspers/cls/secrets"
)

func SetupTLSServer(fp string) {
	log.SetFlags(log.Lshortfile)

	certificate, err := tls.LoadX509KeyPair(secrets.ServerCrtPath, secrets.ServerKeyPath)
	if err != nil {
		log.Printf("Failed to load Keypair: %v", err)
		return
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
	}
	config.MinVersion = tls.VersionTLS13
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
		go HandleConnection(tlsConn, fp)
	}

}

func HandleConnection(conn *tls.Conn, fp string) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	// Listen for marker
	marker, err := filehandler.ListenForMarker(r)
	if err != nil {
		panic(err)
	}

	if marker == "SHARE_FILE_SHARE_FILE" {
		filehandler.HandleFileTransfer(conn, r, fp)
	} else if marker == "LIST_ALL_LIST_ALL" {
		// Log received marker
		log.Printf("Received marker: %v", marker)

		// List all files in server
		files, err := filehandler.GetRemoteFilePaths(fp)
		if err != nil {
			log.Printf("Error retrieving files: %v", err)
		}
		fmt.Printf("Retrieved files: %v\n", files)

		// Respond with "NEXT_ITEM"
		_, err = conn.Write([]byte("NEXT_ITEM\n"))
		if err != nil {
			log.Printf("Error sending response: %v", err)
			return
		}

		for _, file := range files {
			_, err := conn.Write([]byte(file + "\n"))
			if err != nil {
				log.Printf("Error sending file: %v", err)
				return
			}
		}

		// End of list marker
		_, err = conn.Write([]byte("EOL_EOL_EOL_EOL\n"))
		if err != nil {
			log.Printf("Error sending EoL marker: %v", err)
			return
		}

	} else if marker == "GET_FILE_GET_FILE" {
		log.Printf("Received marker: %v", marker)
		sendFile(r, conn, fp)

	}

}

func sendFile(r *bufio.Reader, conn *tls.Conn, fp string) {
	f := arguments.FileData{}

	// Send Share file marker
	_, err := conn.Write([]byte("SHARE_FILE_SHARE_FILE\n"))
	if err != nil {
		log.Printf("Error sending sharefile marker: %v", err)
		return
	}

	// Wait for "NEXT_ITEM"
	if !filehandler.ReadyForNextItem(r) {
		return
	}

	// Get filename
	msg, err := r.ReadString('\n')
	msg = strings.TrimSpace(msg)

    // Wait for "NEXT_ITEM"
    if !filehandler.ReadyForNextItem(r) {
        return
    }

	// Check if file exists
	fullPath := filepath.Join(fp, msg)
	fmt.Printf("Filename: %v\n", fullPath)
	_, err = os.Stat(fullPath)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		// File does not exist
		conn.Write([]byte("File does not exist\n"))
		return
	}

	// Put filedata into 'f'
	fN := filepath.Base(fullPath) // [f]ile[N]ame
	copy(f.Filename[:], fN)
	fE := filepath.Ext(fullPath) // [f]ile[E]xtension
	copy(f.Extension[:], fE)
	fC, err := os.ReadFile(fullPath) // [f]ile[C]ontent
	if err != nil {
		log.Printf("Error reading local file: %v\n", err)
		return
	}
	copy(f.Content, fC)

	// Send filename
	fileName := f.Filename[:]
	_, err = conn.Write(fileName)
	if err != nil {
		log.Printf("Error sending filename: %v", err)
		return
	}
	log.Println("Filename sent, waiting for response...")

	// Wait for "NEXT_ITEM"
	if !filehandler.ReadyForNextItem(r) {
		return
	}

	// Send file extension
	extension := f.Extension[:]
	_, err = conn.Write(extension)
	if err != nil {
		log.Printf("Error sending file extension: %v", err)
		return
	}
	log.Println("File extension sent, waiting for response...")

	// Wait for "NEXT_ITEM"
	if !filehandler.ReadyForNextItem(r) {
		return
	}

	// Send file content
	_, err = conn.Write(fC)
	if err != nil {
		log.Printf("Error sending file content: %v", err)
		return
	}

	log.Println("File content sent, waiting for final confirmation...")

	// Send an EOF marker
	_, err = conn.Write([]byte("EXIT_EOF_EXIT_EOF\n"))
	if err != nil {
		log.Printf("Error sending EOF marker: %v", err)
		return
	}

	// Wait for the final message
	finalMsg, err := r.ReadString('\n')
	if err != nil {
		log.Printf("Error reading final confirmation: %v", err)
		return
	}
	log.Printf("Received from server: %s", strings.TrimSpace(finalMsg))

	log.Println("File transfer complete.")
}
