package main

import (
	"fmt"
	"log"
	"os"

	"github.com/nielsjaspers/cls/internal/arguments"
	"github.com/nielsjaspers/cls/internal/client"
)

func main() {
    clientsideFolder, err := os.UserHomeDir()
    if err != nil {
        log.Fatalf("Error getting user home dir: %v", err)
    }

    os.MkdirAll(clientsideFolder + "/cls-received", os.ModePerm)

	fileContent, args, err := arguments.ExecuteCommand()
	if err != nil {
		fmt.Printf("Error while retrieving file content: %v", err)
		os.Exit(1)
	}
    var argsFixedSize [3]string
    copy(argsFixedSize[:], args)

    client.SetupTLSClient(&fileContent, &argsFixedSize) 

}
