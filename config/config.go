package config

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
