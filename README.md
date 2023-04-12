# distrBoltX
Handwriting based on boltDB distributed KV database, the library will be updated from time to time, suitable for small white scholars entry and distributed advanced
>该库是基于 **io/etcd.bbolt** 驱动打造一个分布式KV库(Bbolt有点类似innodb 完全兼容ACID事务)，新能完全取决于Bbolt的B+tree的顺序写，和MMAP的预随机读，因为是基于硬盘的读写驱动，所以在固态硬盘上运行的性能最佳
# Prepare the dependency library
> go mod tidy 
# Standalone Start to up
> go mod install; main -db-location=shard0.db  -http-addr=127.0.0.2:8080 -config-file=sharding.toml -shard=shard0
# Supporting a simple data sharding,which the server sharding is being accessed
```toml
[[shards]]
name = "shard0"
idx  = 0
address = "127.0.0.2:8080"
replicas = ["127.0.0.22:8080"]

[[shards]]
name = "shard1"
idx  = 1
address = "127.0.0.3:8081"
replicas = ["127.0.0.33:8080"]

[[shards]]
name = "shard2"
idx  = 2
address = "127.0.0.4:8082"
replicas = ["127.0.0.44:8080"]

[[shards]]
name = "shard3"
idx = 3
address = "127.0.0.5:8083"
replicas = ["127.0.0.55:8080"]

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
4. shard 
5. replicas

集群分片采用CRC64 MOD SHARD_COUNT 得到 当前分片，如果有数据写入当前分片，又单协程轮询同步到副本节点，副本节点启动时即刻加载对主节点的写入监听，内部节点采用节点转发的方式避免集群连接过多(参考redis HASHSLOT REDIRECT)

#### DEBUG
When you need to debug locally, here is a more suitable way that you can start using two VS/IDEA boosters for example<br/>
1. Open both editors
2. Enter two startup scripts in the editor, respectively, as follows
```shell

主节点
 --db-location=shard0.db  --http-addr=127.0.0.2:8080  --grpc-addr=127.0.0.2:50030 --config-file=sharding.toml --config=shard0
副本
--db-location=shard0-r.db --http-addr=127.0.0.22:8080 --grpc-addr=127.0.0.2:50030 --config-file=sharding.toml --config=shard0 --replica

```
3. If you only need to start two nodes for testing, you only need to keep one shards shard in sharding.toml as follows:<br/>
```shell
[[shards]]
name = "shard0"
idx  = 0
address = "127.0.0.2:8080"
replicas = ["127.0.0.22:8080"]
```

##### 您可能会问,为什么我在本地可以监听127.0.0.2

#### CHANGELOG_FEATURE
1. 内部连接使用GRPC代替HTTP1.1协议 (done)
2. 取模分片算法采用一致性HASH算发代替用来解决HASH迁移的问题
3. 分片之后数据合并可能会有问题，所有可以参考REDIS HASHTAG 实现HASH聚合 