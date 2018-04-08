package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	//CREATES A SERVER IN LOCALHOST IN THE PORT 8080
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err.Error())
	}

	//CHANNELS FOR INCOMING CONNECTIONS, DEAD CONNECTIONS AND MESSAGES
	conns := make(chan net.Conn)
	dconns := make(chan net.Conn)
	msgs := make(chan string)
	//MAP WITH ALL CONNECTIONS
	aconns := make(map[net.Conn]int)
	//NUMBER OF CONNECTIONS
	i := 0

	//FUNC TO ACCEPT A CONNECTION
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
			}
			conns <- conn
		}
	}()

	for {
		select {
		case conn := <-conns:
			//ADD THE NEW CONNECTION TO THE MAP
			aconns[conn] = i
			//INCREASE THE NUMBER OF CONNECTIONS
			i++
			//ONCE WE HAVE A CONNECTION, WE START READING MESSAGES FROM IT
			go func(conn net.Conn, i int) {
				rd := bufio.NewReader(conn)
				for {
					m, err := rd.ReadString('\n')
					if err != nil {
						break
					}
					msgs <- fmt.Sprintf("Client %v: %v", i, m)
				}
				//DONE READING
				dconns <- conn

			}(conn, i)
		case msg := <-msgs:
			//WE NEED TO BROADCAST THE MSG TO ALL CONNECTIONS
			for conn := range aconns {
				conn.Write([]byte(msg))
			}
		//WHEN A USER IS DISCONECTED
		case dconn := <-dconns:
			log.Printf("Client %v is gone\n", aconns[dconn])
			delete(aconns, dconn)
		}

	}

}
