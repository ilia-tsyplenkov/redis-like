package main

import (
	"juno/server"

	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var globalHash = make(map[string]*server.DataMap)
var host string = "localhost"
var port string = "8000"
var addr string = host + ":" + port
var defalutDbIndex string = "0"

var launchChecker = make(chan string)

func init() {
	var dm server.DataMap
	dm.Init()
	globalHash[defalutDbIndex] = &dm
}

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	go ttlMonitor()
	launchChecker <- defalutDbIndex
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
	prompt := fmt.Sprintf("%s[%s] ", addr, defalutDbIndex)
	dm := globalHash[defalutDbIndex]
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
				_, ok := globalHash[id]
				if !ok {
					globalHash[id] = &server.DataMap{}
					dm = globalHash[id]
					dm.Init()
					launchChecker <- id
				} else {
					dm = globalHash[id]
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

func ttlChecker(dm *server.DataMap) {
	for {
		for _, key := range dm.Keys() {
			ttl, _ := dm.TTL(key)
			if ttl == "-1" {
				continue
			}
			dTTL, err := strconv.Atoi(ttl)
			if err != nil {
				log.Printf("Got unhandled ttl: %q for %q key\n", ttl, key)
				log.Printf("TTL for %q key has been reseted\n", key)
				dm.Persist(key)
				continue
			}
			if int64(dTTL) < time.Now().UTC().Unix() {
				log.Printf("%q key has been removed. TTL expired\n", key)
				dm.Remove(key)
			}
		}
	}
}

func ttlMonitor() {
	for {
		key := <-launchChecker
		go ttlChecker(globalHash[key])
	}
}
