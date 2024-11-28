package main

import (

	"github.com/nielsjaspers/cls/internal/server"
)

func main() {
    serverPath := server.ExecuteRemotePath()
	server.SetupTLSServer(serverPath)
}
