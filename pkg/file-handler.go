package filehandler

import (
	"fmt"
	"os"
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
// returns a string array with paths to the files, and an error if there is one
func GetRemoteFilePaths() ([]string, error) { // Use TLS connection (conn) maybe as parameter ?
    return []string{"",""}, nil
}

// GetRemoteFile requests a single file with remote path p
// Returns a file and an error if there is one
func getRemoteFile(p string) (os.File, error) {
    fmt.Printf("%v", p)
    return *os.NewFile(1, ""), nil
}
