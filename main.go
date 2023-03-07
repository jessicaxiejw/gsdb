package main

import (
	"gsdb/internal/sheet"
	"log"
	"net"
)

// TODO: use a config file for this and move to a package
const (
	HOST = "localhost"
	PORT = "9001"
	TYPE = "tcp"
)

func handleIncomingRequest(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024) // TODO: better buffer allocation
	_, err := conn.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	err = sheet.ParseStatement(string(buffer))
	if err != nil {
		panic(err)
	}
}

func main() {
	listen, err := net.Listen(TYPE, HOST+":"+PORT)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	// TODO: support multiple connections
	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleIncomingRequest(conn)
	}

}
