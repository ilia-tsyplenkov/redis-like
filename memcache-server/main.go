package main

import (
	"redis-like/server"

	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var globalHash = make(map[string]*server.DataMap)
var host string = "localhost"
var port = flag.String("port", "8000", "sever port")
var addr string
var defalutDbIndex string = "0"

var launchChecker = make(chan string)

func init() {
	var dm server.DataMap
	dm.Init()
	dm.DbId = defalutDbIndex
	globalHash[defalutDbIndex] = &dm
}

func main() {
	flag.Parse()
	addr = host + ":" + *port
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

// handleConn handles each c connection.
// it also sends db id through launchChecker channel
// for each new database.
func handleConn(c net.Conn) {
	prompt := fmt.Sprintf("%s[%s] ", addr, defalutDbIndex)
	dm := globalHash[defalutDbIndex]
	input := bufio.NewScanner(c)
	defer c.Close()
	fmt.Fprintf(c, "%s", prompt)
	for input.Scan() {
		cmd, data, err := server.CommandHandler(input.Text())
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
					dm.DbId = id
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

// ttlChecker check ttl for each key
// in endless loop. It resets unhandled ttl to 0.
func ttlChecker(dm *server.DataMap) {
	for {
		for _, key := range dm.Keys() {
			ttl, _ := dm.TTL(key)
			if ttl == "-1" {
				continue
			}
			dTTL, err := strconv.Atoi(ttl)
			if err != nil {
				log.Printf("db %s: Got unhandled ttl %q for %q key\n", dm.DbId, ttl, key)
				log.Printf("db %s: TTL for %q key has been reseted\n", dm.DbId, key)
				dm.Persist(key)
				continue
			}
			if int64(dTTL) < time.Now().UTC().Unix() {
				log.Printf("db %s: %q key has been removed. TTL expired\n", dm.DbId, key)
				dm.Remove(key)
			}
		}
	}
}

// ttlMonitor get db id from launchChecker channel
// and launch ttlChecker got it.
func ttlMonitor() {
	for {
		key := <-launchChecker
		go ttlChecker(globalHash[key])
	}
}
