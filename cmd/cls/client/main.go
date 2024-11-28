package main

import (
	"fmt"
	"github.com/nielsjaspers/cls/internal/arguments"
	"github.com/nielsjaspers/cls/internal/client"
	"os"
)

func main() {
	fileContent, args, err := arguments.ExecuteCommand()
	if err != nil {
		fmt.Printf("Error while retrieving file content: %v", err)
		os.Exit(1)
	}
    var argsFixedSize [3]string
    copy(argsFixedSize[:], args)

	if len(fileContent.Content) > 0 {
		client.SetupTLSClient(&fileContent, &argsFixedSize) 
	} else {
		fmt.Println("No file content processed.")
        os.Exit(1)
	}
}
