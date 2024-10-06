# redis-in-go

An implementation of a Redis server using Go. Inspired by [CodeCrafters Redis challenge](https://app.codecrafters.io/courses/redis/overview). This Redis server is capable of serving basic key-value store commands & allows for creating and reading Redis streams. In addition, this codebase implements RDB persistence, allows for spinning up replica servers and propagating commands, and supports Redis transactions.

## Key-Value Store & Streams

The core key-value store is implemented in the `store` package, providing an interface that allows a client to get, set, delete, and increment keys, in addition to functionality for creating, appending to, and reading from streams. This implementation relies on a Golang mapping from string (key value) to `KeyValue` struct, which is capable of storing either a `string` or `stream` value type. Store operations do rely on read & write locks on mutexes to minimize conflicts. See `store/streams.go` for more information on the streams implementation.

## Running the Server

TODO: include usage - the ./run.sh commands and options + testing with redis-cli
