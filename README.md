# redis-in-go

An implementation of a Redis server using Go. Inspired by [CodeCrafters Redis challenge](https://app.codecrafters.io/courses/redis/overview). This Redis server is capable of serving basic key-value store commands & allows for creating and reading Redis streams. In addition, this codebase implements RDB persistence, allows for spinning up replica servers and propagating commands, and supports Redis transactions.

## Redis Serialization Protocol (RESP)

The `protocol` package contains a partial implementation of [RESP](https://redis.io/docs/latest/develop/reference/protocol-spec/) to allow for standardized communication between the Redis server and some arbitrary number of Redis clients. This extends to decoding (parsing data coming in from clients) & encoding (serializing server responses to send back to clients).

## Key-Value Store & Streams

The core key-value store is implemented in the `store` package, providing an interface that allows a client to get, set, delete, and increment keys, in addition to functionality for creating, appending to, and reading from streams. This implementation relies on a Golang mapping from string (key value) to `KeyValue` struct, which is capable of storing either a `string` or `stream` value type. Store operations do rely on read & write locks on mutexes to minimize conflicts.

See `store/streams.go` & [Redis Streams](https://redis.io/docs/latest/develop/data-types/streams/) for more information on the streams implementation.

## Command Execution

The `commands` package contains primary implementations for Redis commands that are accepted by the server from clients. `commands/socket.go` is responsible for handling a client connection in a dedicated Goroutine, accepting any incoming data, parsing it for commands, and sending it off for execution. This execution occurs in `commands/main.go`, where the command name is read in order to direct the arguments to the appropriate command handler.

## Transactions

Support for transactions is mainly implemented in `commands/transactions.go` in the `commands` package. [Redis Transactions](https://redis.io/docs/latest/develop/interact/transactions/) allow clients to group the execution of multiple Redis commands into a single step.

## RDB Persistence

[Redis Persistence](https://redis.io/docs/latest/operate/oss_and_stack/management/persistence/) is implemented here using RDB files. The `persistence` package is responsible for being able to dump the current state of the key-value store to an RDB file that is saved to disk, and parse an existing RDB file on disk to persist state back into a running server's key-value store.

## Replication

[Redis Replication](https://redis.io/docs/latest/operate/oss_and_stack/management/replication/) is achieved by the `replication` package. Running a replica server causes the replication handshake to be triggered and for mutation commands to be propagated from the master to all replicas, with sync checking possible via the `REPLCONF GETACK` command.

## Running the Server

TODO: include usage - the ./run.sh commands and options + testing with redis-cli
