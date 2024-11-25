package main

import (
	"os"
    "fmt"
	"github.com/nielsjaspers/cls/internal/arguments"
	"github.com/nielsjaspers/cls/internal/client"
)

func main() {
    args := arguments.InitArgs()
	if err := args.Command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
    client.SetupTLSClient()
}
