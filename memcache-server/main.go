package main

import (
	"juno/server"

	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var globalHash = make(map[string]*server.DataMap)
var host string = "localhost"
var port string = "8000"
var addr string = host + ":" + port

func init() {
	var dm server.DataMap
	dm.Init()
	globalHash["0"] = &dm
}

func main() {
	// var dm DataMap
	// dm.Init()
	// globalHash["0"] = dm
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
		go handleConn(conn)
	}
}

func handleConn(c net.Conn) {
	prompt := fmt.Sprintf("%s[0] ", addr)
	dm := globalHash["0"]
	input := bufio.NewScanner(c)
	defer c.Close()
	fmt.Fprintf(c, "%s", prompt)
	for input.Scan() {
		cmd, data, err := server.CommandHandler(input.Text())
		fmt.Fprintf(c, "CommandHandler cmd: %s, data: %s, data_len: %d, err: %s\n", cmd, data, len(data), err)
		if err != nil {
			fmt.Fprintf(c, "%s", prompt)
			continue
		}
		cmd = strings.ToLower(cmd)
		if cmd == "select" {
			if len(data) != 1 {
				fmt.Fprintf(c, "wrong number of arguments for 'select' command\n%s", prompt)
				continue
			} else {
				id := data[0]
				dm, ok := globalHash[id]
				if !ok {
					globalHash[id] = &server.DataMap{}
					dm = globalHash[id]
					dm.Init()
				}
				prompt = fmt.Sprintf("%s[%s] ", addr, id)
				fmt.Fprintf(c, "%s", prompt)
				continue
			}
		}
		result, err := server.DataHandler(dm, cmd, data)
		if err != nil {
			fmt.Fprintf(c, "%s\n%s", err.Error(), prompt)
			continue
		}
		fmt.Fprintf(c, "%s\n%s", result, prompt)
	}
}
