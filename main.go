package main

import (
	"getfile-server/server"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	defer ln.Close()

	server.Serve(ln)
}
