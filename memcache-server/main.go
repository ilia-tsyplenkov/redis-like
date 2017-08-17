package main

import (
	"redis-like/server"

	"flag"
	"log"
	"net"
)

var host string = "localhost"
var port = flag.String("port", "8000", "sever port")
var addr string

func main() {
	flag.Parse()
	addr = host + ":" + *port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go server.HandleConn(conn, addr)
	}
}
