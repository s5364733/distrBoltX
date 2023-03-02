package web

//jack.lei
import (
	"fmt"
	"github.com/s5364733/distrBoltX/db"
	"hash/fnv"
	"io"
	"net/http"
)

// Server contains  HTTP method handler to be used for  the database
type Server struct {
	db         *db.Database
	shardIdx   int
	shardCount int
	addr       map[int]string
}

// NewServer for used to be http endpoint handler
func NewServer(db *db.Database, shardIdx, shardCount int, addr map[int]string) *Server {
	return &Server{
		db:         db,
		shardIdx:   shardIdx,
		shardCount: shardCount,
		addr:       addr,
	}
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.addr[shard] + r.RequestURI
	resp, err := http.Get(url)
	fmt.Fprintf(w, "redirect from config %d to  config %d(%q)\n", s.shardIdx, shard, url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error's redirect request url: %v", err)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

// GetHandler handles get endpoint
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shard := s.GetShard(key)
	value, err := s.db.GetKey(key)

	//当前有可能不是拿的当前分区的数据,例如当前key计算出来的HASH取模分片之后为0 但是请求的是1分区的库，
	//所以这里导航到0分区即可
	if shard != s.shardIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "ShardIdx: %d , cur config :%d ,addr : %q , value = %q ,error=%v ",
		shard,         //KEY 对应的分片路由ID
		s.shardIdx,    //当前分区
		s.addr[shard], //应该拿分区库所在的地址
		value,
		err)
}

// GetShard Calculates the sum hash value  for the key
func (s *Server) GetShard(key string) int {
	n := fnv.New64()
	n.Write([]byte(key))
	return int(n.Sum64() % uint64(s.shardCount))
}

// SetHandler handles set endpoint
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	err := s.db.SetKey(key, []byte(value))
	shard := s.GetShard("key")

	if shard != s.shardIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "Error=%v, shardIdx %d , current config: %d", err, shard, s.shardIdx)
}

// ListenAndServe starts accept request
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}
