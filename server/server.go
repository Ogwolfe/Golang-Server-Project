package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
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

	buffer := make([]byte, 0, 1024) //Capacity 1024, Length 0
	chunk := make([]byte, 1024)     //Length 1024

	for {

		x, err := c.Read(chunk)
		if err != nil {
			return
		}
		buffer = append(buffer, chunk[:x]...)
		clear(chunk)

		fmt.Printf("%s\n", string(buffer))

		req, n, err := parseRequest(buffer)
		if err != nil {
			c.Write([]byte("400: Bad Request"))

			//Remove bad request from buffer
			buffer = buffer[n:]
			continue
		}

		//No complete request found, continue Read() loop
		if n == 0 {
			continue
		}

		//Remove request from buffer
		buffer = buffer[n:]

		//Open file
		file := append([]byte("./files/"), req.Path...)
		contents, err := os.ReadFile(string(file))
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				c.Write([]byte("404: File not found"))
				continue
			}
			return
		}

		f, err := os.Stat(string(file))
		if err != nil {
			return
		}

		s := f.Size()

		response := append([]byte("OK "), []byte(strconv.Itoa(int(s)))...)
		response = append(response, []byte("\n")...)
		response = append(response, contents...)
		_, err = c.Write(response)
		if err != nil {
			return
		}

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
		return req, len(p) + 1, errors.New("Bad Request")
	}

	//We should have a request that's minimum "GET  \n"
	//Set method and get and set path
	req.Method = "GET"
	req.Path = p[4:]

	//len(p) + 1 is the length of the request + '\n' in the buffer
	return req, len(p) + 1, nil
}
