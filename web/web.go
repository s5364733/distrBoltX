package web

//jack.lei
import (
	"fmt"
	"github.com/s5364733/distrBoltX/config"
	"github.com/s5364733/distrBoltX/db"
	"io"
	"net/http"
)

// Server contains  HTTP method handler to be used for  the database
type Server struct {
	db *db.Database
	//shardIdx   int
	//shardCount int
	//addr       map[int]string
	shards *config.Shards
}

// NewServer for used to be http endpoint handler
func NewServer(db *db.Database, s *config.Shards) *Server {
	return &Server{
		db:     db,
		shards: s,
	}
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.shards.Addrs[shard] + r.RequestURI
	resp, err := http.Get(url)
	fmt.Fprintf(w, "redirect from config %d to  config %d(%q)\n", s.shards.CurIdx, shard, url)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error's redirect request url: %v", err)
		return
	}
	//调用者关闭BODY
	defer resp.Body.Close()
	//写入响应体stream
	io.Copy(w, resp.Body)
}

// GetHandler handles get endpoint
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	shard := s.shards.Index(key)
	value, err := s.db.GetKey(key)

	//当前有可能不是拿的当前分区的数据,例如当前key计算出来的HASH取模分片之后为0 但是请求的是1分区的库，
	//所以这里导航到0分区即可
	if shard != s.shards.CurIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "ShardIdx: %d , cur config :%d ,addr : %q , value = %q ,error=%v ",
		shard,                 //KEY 对应的分片路由ID
		s.shards.CurIdx,       //当前分区
		s.shards.Addrs[shard], //应该拿分区库所在的地址
		value,
		err)
}

// SetHandler handles set endpoint
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	err := s.db.SetKey(key, []byte(value))
	shard := s.shards.Index("key")

	if shard != s.shards.CurIdx {
		s.redirect(shard, w, r)
		return
	}

	fmt.Fprintf(w, "Error=%v, shardIdx %d , current config: %d", err, shard, s.shards.CurIdx)
}

// ListenAndServe starts accept request
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}

// DeleteExtraKeyHandler deletes keys that don't belong to current shard
func (s *Server) DeleteExtraKeyHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error  = %v", s.db.DeleteExtraKeys(func(key string) bool {
		return s.shards.CurIdx != s.shards.Index(key)
	}))
}
