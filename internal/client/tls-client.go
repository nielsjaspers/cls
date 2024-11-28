package client

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/nielsjaspers/cls/internal/arguments"
	"github.com/nielsjaspers/cls/secrets"
)

func SetupTLSClient(f *arguments.FileData, args *[3]string) {
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

	r := bufio.NewReader(conn)

	// Process commands based on args[0]
	switch args[0] {
	case "SHARE_FILE_SHARE_FILE":
		SendFile(f, args, r, conn)
	case "LIST_ALL_LIST_ALL":
		fileList, err := getList(args, r, conn)
		if err != nil {
			log.Fatalf("Error while listing files: %v", err)
		}
		log.Printf("Retrieved files: %v", fileList)
	case "GET_FILE_GET_FILE":
	default:
		log.Println("Invalid command")
		os.Exit(1)
	}
}

func waitForNextItem(r *bufio.Reader) bool {
	msg, err := r.ReadString('\n')
	if err != nil {
		log.Printf("Error reading server response: %v", err)
		return false
	}

	msg = strings.TrimSpace(msg)
	if msg != "NEXT_ITEM" {
		log.Printf("Unexpected server response: %s", msg)
		return false
	}

	log.Println("Received NEXT_ITEM, proceeding to the next step.")
	return true
}

func getList(args *[3]string, r *bufio.Reader, conn *tls.Conn) ([]string, error) {
	var remoteFiles []string

	// Send get list marker
	_, err := conn.Write([]byte(args[0] + "\n"))
	if err != nil {
		return nil, fmt.Errorf("error sending list all marker: %v", err)
	}

	// Wait for "NEXT_ITEM"
	if !waitForNextItem(r) {
		return nil, fmt.Errorf("did not receive expected NEXT_ITEM response")
	}

	// Read the file list
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return nil, fmt.Errorf("unexpected EOF while reading server response")
			}
			return nil, fmt.Errorf("error reading server response: %v", err)
		}

		// Clean up the message (remove trailing spaces/newlines)
		msg = strings.TrimSpace(msg)

		// Check if the End of List marker has been reached
		if msg == "EOL_EOL_EOL_EOL" {
			break
		}

		// Append the file name to the list
		if msg != "" { // Avoid appending empty lines
			remoteFiles = append(remoteFiles, msg)
		}
	}

	return remoteFiles, nil
}

func SendFile(f *arguments.FileData, args *[3]string, r *bufio.Reader, conn *tls.Conn) {
	// Send Share file marker
	_, err := conn.Write([]byte(args[0]))
	if err != nil {
		log.Printf("Error sending sharefile marker: %v", err)
		return
	}

	// Wait for "NEXT_ITEM"
	if !waitForNextItem(r) {
		return
	}

	// Send filename
	fileName := f.Filename[:]
	_, err = conn.Write(fileName)
	if err != nil {
		log.Printf("Error sending filename: %v", err)
		return
	}
	log.Println("Filename sent, waiting for response...")

	// Wait for "NEXT_ITEM"
	if !waitForNextItem(r) {
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
	if !waitForNextItem(r) {
		return
	}

	// Send file content
	_, err = conn.Write(f.Content)
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

