# redis-in-go

An implementation of a Redis server using Go. Inspired by [CodeCrafters Redis challenge](https://app.codecrafters.io/courses/redis/overview). This Redis server is capable of serving basic key-value store commands & allows for creating and reading Redis streams. In addition, this codebase implements RDB persistence, allows for spinning up replica servers and propagating commands, and supports Redis transactions.

## Redis Serialization Protocol (RESP)

The `protocol` package contains a partial implementation of [RESP](https://redis.io/docs/latest/develop/reference/protocol-spec/) to allow for standardized communication between the Redis server and some arbitrary number of Redis clients. This extends to decoding (parsing data coming in from clients) & encoding (serializing server responses to send back to clients).

## Key-Value Store & Streams

The core key-value store is implemented in the `store` package, providing an interface that allows a client to get, set, delete, and increment keys, in addition to functionality for creating, appending to, and reading from streams. This implementation relies on a Golang mapping from string (key value) to `KeyValue` struct, which is capable of storing either a `string` or `stream` value type. Store operations do rely on read & write locks on mutexes to minimize conflicts. See `store/streams.go` for more information on the streams implementation.

## Running the Server

TODO: include usage - the ./run.sh commands and options + testing with redis-cli
