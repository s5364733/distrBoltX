package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/s5364733/distrBoltX/config"
	"github.com/s5364733/distrBoltX/db"
	"github.com/s5364733/distrBoltX/web"
	"log"
	"net/http"
)

var (
	dbLocation = flag.String("db-location", "", "the path to the bolt database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP post and host")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard      = flag.String("config", "", "the name of config for the data")
)

func parseFlag() {
	flag.Parse()
	if *dbLocation == "" {
		log.Fatalf("Must be Provide  db-location")
	}
	if *shard == "" {
		log.Fatalf("Must be Provide  config")
	}
}

func main() {
	// Open the XXXX.db data file in your current directory.
	// It will be created if it doesn't exist.
	parseFlag()

	var c config.Config
	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("Toml.DecodeFile(%q): %v", *configFile, err)
	}

	//var shardCount int
	var shardIdx = -1
	var addrs = make(map[int]string)
	//扫描所有shard,
	for _, s := range c.Shards {
		addrs[s.Idx] = s.Address
		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatalf("Shard %q was not found", *shard)
	}

	log.Printf("Shard count is %d current config :%d", len(c.Shards), shardIdx)
	fmt.Printf("%#v", &c)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q) : %v", *dbLocation, err)
	}

	defer close()
	//shard0 shard1 shard2 分别放在三个数据库
	srv := web.NewServer(db, shardIdx, len(c.Shards), addrs)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/get", srv.GetHandler)

	srv.ListenAndServe(*httpAddr)
}
