# distrBoltX
Handwriting based on boltDB distributed KV database, the library will be updated from time to time, suitable for small white scholars entry and distributed advanced
>该库是基于 **io/etcd.bbolt** 驱动打造一个分布式KV库(Bbolt有点类似innodb 完全兼容ACID事务)，新能完全取决于Bbolt的B+tree的顺序写，和MMAP的预随机读，因为是基于硬盘的读写驱动，所以在固态硬盘上运行的性能最佳
# Prepare the dependency library
> go mod tidy 
# Standalone Start to up
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

[[shards]]
name = "shard3"
idx  = 3
address = "localhost:8083"

```
# Middleware dependency
>bbolt 
> https://github.com/etcd-io

bbolt is a fork of Ben Johnson's Bolt key/value store. The purpose of this fork is to provide the Go community with an active maintenance and development target for Bolt; the goal is improved reliability and stability. bbolt includes bug fixes, performance enhancements, and features not found in Bolt while preserving backwards compatibility with the Bolt API.

Bolt is a pure Go key/value store inspired by Howard Chu's LMDB project. The goal of the project is to provide a simple, fast, and reliable database for projects that don't require a full database server such as Postgres or MySQL.

Since Bolt is meant to be used as such a low-level piece of functionality, simplicity is key. The API will be small and only focus on getting values and setting values. That's it.

DistrBoltX is secondary developed based on bbolt, adding distributed fragmentation high availability data security and other scenarios

There will be a lot of optimization details in the future, so stay tuned

# Distributed startup
1. ./populate.sh
2.  检查toml配置文件是否对应服务器完整
3. ./luanch.sh

当你看到:<br/>
![](img/c18e797d7c4525afd03a7ff1e85e014.png)
说明你此时已经启动了四个端口监听四个分片库 ,You know ?

Core module :
1. 数据分片 
2. 读写基准测试
3. 多机备份
4. reshard 
5. replication(分区)
