# distrBoltX
Handwriting based on boltDB distributed KV database, the library will be updated from time to time, suitable for small white scholars entry and distributed advanced
# Prepare the dependency library
> go mod tidy 
# Start to up
> go mod install; main -db-location=distBoltX.db
# Supporting a simple data sharding,which the server sharding is being accessed
```toml
[[shards]]
name = "shard0"
idx  = 0
address = "localhost:8080"

[[shards]]
name = "shard1"
idx  = 1
address = "localhost:8081"

[[shards]]
name = "shard2"
idx  = 2
address = "localhost:8082"
```
# Middleware dependency
>bbolt 
> https://github.com/etcd-io

bbolt is a fork of Ben Johnson's Bolt key/value store. The purpose of this fork is to provide the Go community with an active maintenance and development target for Bolt; the goal is improved reliability and stability. bbolt includes bug fixes, performance enhancements, and features not found in Bolt while preserving backwards compatibility with the Bolt API.

Bolt is a pure Go key/value store inspired by Howard Chu's LMDB project. The goal of the project is to provide a simple, fast, and reliable database for projects that don't require a full database server such as Postgres or MySQL.

Since Bolt is meant to be used as such a low-level piece of functionality, simplicity is key. The API will be small and only focus on getting values and setting values. That's it.

DistrBoltX is secondary developed based on bbolt, adding distributed fragmentation high availability data security and other scenarios

There will be a lot of optimization details in the future, so stay tuned