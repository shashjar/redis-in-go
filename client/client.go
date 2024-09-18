package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const SERVER_TYPE = "tcp"
const HOSTNAME = "0.0.0.0"
const PORT = "6379"

func openSocketConnection() net.Conn {
	address := HOSTNAME + ":" + PORT
	conn, err := net.Dial(SERVER_TYPE, address)
	if err != nil {
		fmt.Println("Failed to connect client to Redis server:", err.Error())
		os.Exit(1)
	}
	return conn
}

func main() {
	args := os.Args[1:]

	conn := openSocketConnection()
	defer conn.Close()

	redisCommand := strings.Join(args, " ")
	_, err := fmt.Fprintf(conn, "%s\n", redisCommand)
	if err != nil {
		fmt.Println("Error sending data to server:", err.Error())
		os.Exit(1)
	}

	buf := make([]byte, 128)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Failed to read response from Redis server:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Response from Redis server:", string(buf))
}
