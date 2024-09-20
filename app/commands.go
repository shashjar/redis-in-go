package main

import (
	"log"
	"net"
	"os"
)

func ping(conn net.Conn) {
	_, err := conn.Write([]byte(toSimpleString("PONG")))
	if err != nil {
		log.Println("Error writing PONG to connection: ", err.Error())
		os.Exit(1)
	}
}

func echo(conn net.Conn, command []string) {
	if len(command) <= 1 {
		_, err := conn.Write([]byte(toSimpleError("ERR wrong number of arguments for 'echo' command")))
		if err != nil {
			log.Println("Error writing to connection:", err.Error())
			os.Exit(1)
		}
	}

	_, err := conn.Write([]byte(toBulkString(command[1])))
	if err != nil {
		log.Println("Error executing ECHO on connection:", err.Error())
		os.Exit(1)
	}
}

func commandDocs(conn net.Conn) {
	_, err := conn.Write([]byte("*0\r\n"))
	if err != nil {
		log.Println("Error executing COMMAND DOCS on connection:", err.Error())
		os.Exit(1)
	}
}
