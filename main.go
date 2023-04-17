package main

import (
	"flag"
	"fmt"
	"github.com/s5364733/distrBoltX/api"
	"github.com/s5364733/distrBoltX/config"
	"github.com/s5364733/distrBoltX/internal/db"
	"github.com/s5364733/distrBoltX/internal/rpc/serv"
	"github.com/s5364733/distrBoltX/internal/web"
	"github.com/s5364733/distrBoltX/pkg/replication"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

var (
	dbLocation = flag.String("db-location", "", "the path to the bolt database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP post and host")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard      = flag.String("config", "", "the name of config for the data")
	replica    = flag.Bool("replica", false, "Whether or not run as a read-only replica")
	grpcAddr   = flag.String("grpc-addr", "127.0.0.1:50030", "grpc's inner port register ")
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
	log.Printf("Shard count is %d current shard :%d cur config %#v:", len(c.Shards), shards.CurIdx, &c)

	db, close, err := db.NewDatabase(*dbLocation, *replica)
	if err != nil {
		log.Fatalf("NewDatabase(%q) : %v", *dbLocation, err)
	}

	defer close()
	//如果当前是副本
	//shard0 shard1 shard2 分别放在三个数据库
	srv := registerHttpFuncHandler(db, shards)

	if *replica {
		//拿到当前主分片节点 fetch current master node address
		//name = "shard0"
		//idx  = 0
		//address = "127.0.0.2:8080"
		//replicas = ["127.0.0.22:8080"]
		//这里就是 shard0 == 	127.0.0.2:8080
		leaderAddr, ok := shards.Addrs[shards.CurIdx]
		if !ok {
			log.Fatalf("Could not find address for leader for shard %d", shards.CurIdx)
		}
		//启动一个协程去轮询
		go replication.ClientGrpcLoop(db, leaderAddr)
	} else { //GRPC 端口注册
		fmt.Printf("execute init for grpc register %#v register node ip addr : %q", srv, *grpcAddr)
		go registerGrpcPort(srv, *grpcAddr)
	}

	//开启主节点同步端口
	srv.ListenAndServe(*httpAddr)

}

func registerHttpFuncHandler(db *db.Database, shards *config.Shards) *web.Server {
	srv := web.NewServer(db, shards)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/purge", srv.DeleteExtraKeyHandler)
	http.HandleFunc("/next-replication-key", srv.GetNextKeyForReplication)
	http.HandleFunc("/delete-replication-key", srv.DeleteReplicationKey)
	return srv
}

func registerGrpcPort(server *web.Server, grpcAddr string) {
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	api.RegisterAckSyncDialerServer(s, serv.NewAckSyncDialerService(server))
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		s.GracefulStop()
	}
}
