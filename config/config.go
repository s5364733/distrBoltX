package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"hash/fnv"
)

// Shard describes a config that holds the key
// Each key has unique the set of key
type Shard struct {
	Name    string
	Idx     int
	Address string
}

// Config describes the sharding config
type Config struct {
	Shards []Shard
}

// Shards represent an easier to use representation of
// the sharding  config : the shards count, current index and
// the addresses of all other  shards too
type Shards struct {
	Count  int            //分片总数
	CurIdx int            //当前分片索引
	Addrs  map[int]string //分片索引对地址MAPPING
}

// ParseConfig parses the config and returns it ~~~~~~~~~~~~~~~~~~~~
func ParseConfig(filename string) (Config, error) {
	var c Config
	if _, err := toml.DecodeFile(filename, &c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// ParseShards converts and verifies the list of shards
// specified in the config into a form that can be used
// for routing
// 稍微抽象一下路由解析方法
func ParseShards(shards []Shard, curShardName string) (*Shards, error) {
	shardCount := len(shards)
	shardIdx := -1
	addrs := make(map[int]string) //[Idx,addr]
	//遍历所有分片
	for _, s := range shards {
		if _, ok := addrs[s.Idx]; ok {
			//ok == true 说明有值，反之说明没有值
			return nil, fmt.Errorf("duplicate shard index:%d", s.Idx)
		}
		addrs[s.Idx] = s.Address
		if s.Name == curShardName {
			shardIdx = s.Idx //拿到当前shardIdx
		}
	}
	for i := 0; i < shardCount; i++ {
		if _, ok := addrs[i]; !ok {
			return nil, fmt.Errorf("shard %d is not found", i)
		}
	}
	if shardIdx < 0 {
		return nil, fmt.Errorf("shard %q was not found", curShardName)
	}
	//拿到所有分片和地址
	//拿到当前分片Idx
	return &Shards{
		Addrs:  addrs,
		Count:  shardCount,
		CurIdx: shardIdx,
	}, nil
}

// Index return the shard number for the corresponding key
func (s *Shards) Index(key string) int {
	h := fnv.New64()
	h.Write([]byte(key))
	return int(h.Sum64() % uint64(s.Count))
}

//// GetShard Calculates the sum hash value  for the key
//func (s *Server) GetShard(key string) int {
//	n := fnv.New64()
//	n.Write([]byte(key))
//	return int(n.Sum64() % uint64(s.shards.Count))
//}
