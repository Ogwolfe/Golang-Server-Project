package server

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

type Request struct {
	Method string
	Path   string
}

func StartServer() {
	//Setup Listener
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	buffer := make([]byte, 0, 1024)
	chunk := make([]byte, 1024)

	len := 0

	for {

		x, err := c.Read(chunk)
		if err != nil {
			return
		}
		buffer = append(buffer, chunk[:x]...)
		len += x

		req, n, err := parseRequest(buffer[:len])
		if err != nil {
			c.Write([]byte("400: Bad Request"))
		}

		if n == 0 {
			continue
		}

		_, err = c.Write(append([]byte(req.Method), []byte(req.Path)...))
		if err != nil {
			return
		}

		//Remove request from buffer
		buffer = buffer[n:]
		len -= n
	}

}

func parseRequest(b []byte) (r Request, n int, e error) {
	req := Request{"", ""}
	//check if contains '\n'
	s := string(b)

	if !strings.ContainsAny(s, "\n") {
		return req, 0, nil
	}

	x := strings.SplitN(s, "\n", 2)
	if x == nil {
		return req, 0, errors.New("Error splitting string")
	}

	p := x[0]
	//Check if request is min length then check if it is a proper "GET " request
	if len(p) < 4 || strings.Compare(p[:4], "GET ") != 0 {
		fmt.Printf("p: %q\n", p)
		fmt.Printf("first 4: %q\n", p[:4])
		fmt.Printf("first 4 bytes: %v\n", []byte(p[:4]))
		return req, 0, errors.New("Bad Request")
	}

	//We should have a request that's minimum "GET  \n"
	//Set method and get and set path
	req.Method = "GET"
	req.Path = p[4:]

	//len(p) + 1 is the length of the request + '\n' in the buffer
	return req, len(p) + 1, nil
}
