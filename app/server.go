package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Running server...")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading data into buffer: ", err.Error())
		os.Exit(1)
	}
	log.Printf("Read command: \n%s\n", buf)

	_, err = conn.Write([]byte("+PONG\r\n"))
	if err != nil {
		fmt.Println("Error writing PONG to connection: ", err.Error())
		os.Exit(1)
	}
}
