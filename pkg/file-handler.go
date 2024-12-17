package filehandler

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileUpload reads file from local path p
// returns a byte array with file content, and an error if there is one
func FileUpload(p string) ([]byte, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetRemoteFilePaths requests all remote paths of files currently on the server
// returns a string array with paths to the files, or an error if there is one
func GetRemoteFilePaths(path string) ([]string, error) {
	var allFiles []string

	files, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return []string{"EOF"}, err
	}

	for _, file := range files {
		if !file.IsDir() { // Check if it's not a directory
			allFiles = append(allFiles, file.Name())
		}
	}
	return allFiles, nil
}

// HandleFileTransfer transfers a file over a network connection, using conn for communication, r to read data, and fp as the file path prefix for saving the received file.
func HandleFileTransfer(conn *tls.Conn, r *bufio.Reader, fp string) {
	// Respond with "NEXT_ITEM"
	_, err := conn.Write([]byte("NEXT_ITEM\n"))
	if err != nil {
		log.Printf("Error sending response: %v", err)
		return
	}

	// Listen for filename (max 255 bytes)
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

	var filePath string
	if strings.HasPrefix(fp, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Error getting home directory: %v", err)
			return
		}
		fp = filepath.Join(homeDir, fp[1:])
	}
	if fp == "" {
		filePath = fmt.Sprintf("%s", filename)
	} else {
		filePath = fmt.Sprintf("%s/%s", fp, filename)
	}

	// Open the file for writing
	file, err := os.Create(filePath)
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

func ReadyForNextItem(r *bufio.Reader) bool {
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

func ListenForMarker(r *bufio.Reader) (string, error) {
	markerBuf := make([]byte, 32)
	m, err := r.Read(markerBuf)
	if err != nil {
		return "", fmt.Errorf("Error reading marker: %v\n", err)
	}
	marker := string(bytes.Trim(markerBuf[:m], "\x00\n"))
	return marker, nil
}
