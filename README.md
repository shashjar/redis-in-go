[![Go 1.22](https://img.shields.io/badge/go-1.22-9cf.svg)](https://golang.org/dl/)

# redis-in-go

An implementation of a Redis server using Go. Inspired by the [CodeCrafters Redis challenge](https://app.codecrafters.io/courses/redis/overview). This Redis server is capable of serving basic key-value store commands & allows for creating and reading Redis streams. In addition, this codebase implements RDB persistence, allows for spinning up replica servers and propagating commands, and supports Redis transactions.

![Redis Logo](./docs/assets/redis-logo.png)

## Redis Serialization Protocol (RESP)

The `protocol` package contains a partial implementation of [RESP](https://redis.io/docs/latest/develop/reference/protocol-spec/) to allow for standardized communication between the Redis server and some arbitrary number of Redis clients. This extends to decoding (parsing data coming in from clients) & encoding (serializing server responses to send back to clients).

## Key-Value Store: Strings, Streams, & Lists

The core key-value store is implemented in the `store` package, providing an interface that allows a client to:
- Get, set, delete, and increment string keys
- Create, append to, and read from streams
- Create, add elements to, list elements of, remove elements from, and block retrieval on lists

This key-value store implementation relies on a Golang mapping from string (key value) to `KeyValue` struct, which is capable of storing either a `string`, `stream`, or `list` value type. Store operations rely on read & write locks on mutexes to minimize conflicts.

See `store/streams.go` & [Redis Streams](https://redis.io/docs/latest/develop/data-types/streams/) for more information on the streams implementation.

See `store/lists.go` & [Redis Lists](https://redis.io/docs/latest/develop/data-types/lists/) for more information on the lists implementation.

## Command Execution

The `commands` package contains primary implementations for Redis commands that are accepted by the server from clients. `commands/socket.go` is responsible for handling a client connection in a dedicated Goroutine, accepting any incoming data, parsing it for commands, and sending it off for execution. This execution occurs in `commands/main.go`, where the command name is read in order to direct the arguments to the appropriate command handler.

## Transactions

Support for transactions is mainly implemented in `commands/transactions.go` in the `commands` package. [Redis Transactions](https://redis.io/docs/latest/develop/interact/transactions/) allow clients to group the execution of multiple Redis commands into a single step.

## RDB Persistence

[Redis Persistence](https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/) is implemented here using RDB files. The `persistence` package is responsible for being able to dump the current state of the key-value store to an RDB file that is saved to disk, and parse an existing RDB file on disk to persist state back into a running server's key-value store.

## Replication

[Redis Replication](https://redis.io/docs/latest/operate/oss_and_stack/management/replication/) is achieved by the `replication` package. Running a replica server causes the replication handshake to be triggered and for mutation commands to be propagated from the master to all replicas, with sync checking possible via the `REPLCONF GETACK` command.

## Running the Server

The `./run.sh` script is used to run the server. The program entrypoint accepts the following command-line arguments: `port`, `replicaof`, `dir`, & `dbfilename`. `port` is the port number on which to run the Redis server. `replicaof` should be provided with a string value "<MASTER_HOST> <MASTER_PORT>" indicating the master which this server is replicating. `dir` & `dbfilename` together provide the path of an RDB file from which to persist data on startup, and to save data to when dumping an RDB file to disk (see `persistence` above for more details). The below are valid ways to run a Redis server:

`./run.sh`

`./run.sh --port 6380`

`./run.sh --port 6380 --replicaof "localhost 6379"`

`./run.sh --dir /redis-data --dbfilename dump.rdb`

To spin up an official Redis server to compare/test responses & commands, `redis-server` can be used. To spin up one or more clients to connect to the server and run commands, `redis-cli` can be used.
