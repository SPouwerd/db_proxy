package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":5432")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	log.Println("Listening to TCP requests on port", ln.Addr().String())

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Print("Received from ", conn.RemoteAddr(), " \n", buf[:n], "\n")
	}
}
