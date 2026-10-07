Author: Jason Wolfe

This is a TCP server written in Go implementing a simple protocol for downloading files from a server
The goal is for me to learn Go and be productive


Protocol: <METHOD> <PATH>/n

METHOD must be "GET"
PATH specifies the name of a file to receive

Request must be terminated with a newline '\n' character and have a space ' ' between METHOD and PATH

Reading strategy:

handleConnection() will maintain a dynamic slice to act as a buffer for read in bytes.
After each Read() it calls parseRequest(buffer). 
If a complete request has been read into the buffer then it returns a filled out Request struct.
