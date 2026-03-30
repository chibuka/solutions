package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error: %v\n", err)
		}

		r := bufio.NewReader(conn)
		line, _, _ := r.ReadLine()

		req := strings.Fields(string(line))
		path := req[1]

		if path == "/" {
			conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
		} else {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
		}

	}
}
