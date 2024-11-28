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

// GetRemoteFile requests a single file with remote path p
// Returns a file and an error if there is one
func getRemoteFile(p string) (os.File, error) {
	fmt.Printf("%v", p)
	return *os.NewFile(1, ""), nil
}
