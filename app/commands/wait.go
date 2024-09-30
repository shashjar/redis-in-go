package commands

import (
	"net"
	"strconv"
	"time"

	"github.com/shashjar/redis-in-go/app/protocol"
	"github.com/shashjar/redis-in-go/app/replication"
)

// WAIT command
func wait(conn net.Conn, command []string) {
	if len(command) != 3 {
		write(conn, protocol.ToSimpleError("ERR wrong number of arguments for WAIT command"))
		return
	}

	numReplicas, err := strconv.Atoi(command[1])
	if err != nil || numReplicas > len(replication.SERVER_CONFIG.Replicas) {
		write(conn, protocol.ToSimpleError("ERR invalid number of replicas for WAIT command"))
		return
	}

	timeout, err := strconv.Atoi(command[2])
	if err != nil {
		write(conn, protocol.ToSimpleError("ERR invalid timeout for WAIT command"))
		return
	}

	numReplicasAcknowledged := waitForReplicaAcknowledgement(numReplicas, timeout)
	write(conn, protocol.ToInteger(numReplicasAcknowledged))
}

func waitForReplicaAcknowledgement(numReplicas int, timeout int) int {
	var numReplicasAcknowledged int
	acknowledgedReplicas := make(map[string]struct{})

	var timerChannel <-chan time.Time
	if timeout > 0 {
		timer := time.NewTimer(time.Duration(timeout) * time.Millisecond)
		timerChannel = timer.C
	}

	replication.SendGetAckToReplicas()
	defer replication.AddGetAckBytesToMasterReplicationOffset()

	for {
		select {
		case <-timerChannel:
			return numReplicasAcknowledged
		default:
			for _, replica := range replication.SERVER_CONFIG.Replicas {
				_, ok := acknowledgedReplicas[replica.ID()]
				if !ok {
					if replica.LastAcknowledgedReplicationOffset == replication.SERVER_CONFIG.MasterReplicationOffset {
						acknowledgedReplicas[replica.ID()] = struct{}{}
						numReplicasAcknowledged += 1

						if numReplicasAcknowledged >= numReplicas {
							return numReplicasAcknowledged
						}
					}
				}
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}
