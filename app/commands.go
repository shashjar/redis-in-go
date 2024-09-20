package main

import (
	"net"
)

// COMMAND DOCS command
func commandDocs(conn net.Conn) {
	write(conn, "*0\r\n")
}

// PING command
func ping(conn net.Conn) {
	write(conn, toSimpleString("PONG"))
}

// ECHO command
func echo(conn net.Conn, command []string) {
	if len(command) <= 1 {
		write(conn, toSimpleError("ERR wrong number of arguments for 'echo' command"))
		return
	}

	write(conn, toBulkString(command[1]))
}

// GET command
func get(conn net.Conn, command []string) {
	if len(command) != 2 {
		write(conn, toSimpleError("ERR wrong number of arguments for 'get' command"))
		return
	}

	val, ok := REDIS_STORE[command[1]]
	if ok {
		write(conn, toBulkString(val))
	} else {
		write(conn, toNullBulkString())
	}
}

// SET command
func set(conn net.Conn, command []string) {
	if len(command) != 3 {
		write(conn, toSimpleError(("ERR wrong number of arguments for 'set' command")))
		return
	}

	REDIS_STORE[command[1]] = command[2]
	write(conn, toSimpleString("OK"))
}
