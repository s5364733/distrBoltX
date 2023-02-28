package web

//jack.lei
import (
	"fmt"
	db "github.com/s5364733/distrBoltX"
	"net/http"
)

// Server contains  HTTP method handler to be used for  the database
type Server struct {
	db *db.Database
}

func NewServer(db *db.Database) *Server {
	return &Server{
		db: db,
	}

}

// GetHandler handles get endpoint
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value, err := s.db.GetKey(key)
	fmt.Fprintf(w, "value = %q ,error=%v", value, err)
}

// SetHandler handles set endpoint
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "value = %q ,error=%v", value, err)
}

// ListenAndServe starts accept request
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}
