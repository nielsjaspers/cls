package main

import (
	"fmt"
	"os"

	"github.com/nielsjaspers/cls/internal/server"
)

func main() {
	args, _ := server.InitRemotePath()
	if err := args.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	server.SetupTLSServer()
}
