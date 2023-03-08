package main

import (
	"gsdb/internal/sheet"
	"gsdb/internal/sql"
	"io/ioutil"
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

	cred, err := ioutil.ReadFile("./credential.json")
	if err != nil {
		log.Fatal(err) // TODO: wrap error
	}
	client, err := sheet.New(cred) // TODO: add config to accept cred from file or from env
	if err != nil {
		log.Fatal(err) // TODO: wrap error
	}

	err = sql.NewPostgreSQL(client).Execute(string(buffer))
	if err != nil {
		log.Fatal(err) // TODO: wrap error
	}
}

func main() {
	cred, err := ioutil.ReadFile("./credential.json")
	if err != nil {
		log.Fatal(err) // TODO: wrap error
	}
	client, err := sheet.New(cred) // TODO: add config to accept cred from file or from env
	if err != nil {
		log.Fatal(err) // TODO: wrap error
	}
	sql.NewPostgreSQL(client)

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
