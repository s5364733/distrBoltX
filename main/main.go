package main

import (
	"flag"
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

	c, err := config.ParseConfig(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config %q: %v", *configFile, err)
	}

	shards, err := config.ParseShards(c.Shards, *shard)
	if err != nil {
		log.Fatalf("Error parsing shards config :%v", err)
	}
	log.Printf("Shard count is %d current config :%d cur config %#v:", len(c.Shards), shards.CurIdx, &c)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabase(%q) : %v", *dbLocation, err)
	}

	defer close()
	//shard0 shard1 shard2 分别放在三个数据库
	srv := web.NewServer(db, shards)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/purge", srv.DeleteExtraKeyHandler)

	srv.ListenAndServe(*httpAddr)
}

////var shardCount int
//var shardIdx = -1
//var addrs = make(map[int]string)
////扫描所有shard,
//for _, s := range c.Shards {
//	addrs[s.Idx] = s.Address
//	if s.Name == *shard {
//		shardIdx = s.Idx
//	}
//}
//
//if shardIdx < 0 {
//	log.Fatalf("Shard %q was not found", *shard)
//}
