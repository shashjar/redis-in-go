package commands

import "net"

// TODO: returns an empty array to get redis-cli to initialize properly
// COMMAND DOCS command
func commandDocs(conn net.Conn) {
	write(conn, "*0\r\n")
}
