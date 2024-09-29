package commands

import (
	"log"
	"net"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/store"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	log.Println("Handling connection")
	for {
		command, err := readCommand(conn)
		log.Println("Read command:", command, err)
		if err != nil {
			log.Println("Error reading command from connection", err.Error())
			return
		}

		if len(command) > 0 {
			executeCommand(command, conn)
		}
	}
}

// Writes the provided message to the provided connection, but only if this server is not a replica, since
// replicas should not send responses back to the master after receiving propagated commands.
func write(conn net.Conn, message string) {
	if !store.SERVER_CONFIG.IsReplica {
		_, err := conn.Write([]byte(message))
		if err != nil {
			log.Println("Error writing to connection:", err.Error())
			conn.Close()
		}
	}
}

func readIntoBuffer(conn net.Conn) (int, []byte, error) {
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		return n, nil, err
	}

	if n == 0 {
		return 0, []byte{}, nil
	}

	buf = buf[:n]

	return n, buf, nil
}

func readCommand(conn net.Conn) ([]string, error) {
	n, buf, err := readIntoBuffer(conn)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return []string{}, nil
	}

	cmd, err := protocol.ParseCommand(buf)
	if err != nil {
		return nil, err
	}

	return cmd, nil
}