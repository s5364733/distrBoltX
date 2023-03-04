package web

import (
	"fmt"
	"github.com/s5364733/distrBoltX/db"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func createShardDb(t *testing.T, idx int) *db.Database {
	t.Helper()
	f, err := os.CreateTemp(os.TempDir(), fmt.Sprintf("db%d", idx))
	if err != nil {
		t.Fatalf("Could not create temp file %v", err)
	}
	f.Close()
	name := f.Name()
	t.Cleanup(func() {
		os.Remove(name)
	})

	database, closeFunc, err := db.NewDatabase(name)
	if err != nil {
		t.Fatalf("Could not create new db err :%v", err)
	}

	t.Cleanup(func() {
		closeFunc()
	})

	return database
}

func TestNewServer(t *testing.T) {
	var ts1GetHandler, ts1SetHandler func(w http.ResponseWriter, r *http.Request)
	var ts2GetHandler, ts2SetHandler func(w http.ResponseWriter, r *http.Request)

	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			ts1SetHandler(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			ts1GetHandler(w, r)
		}
	}))

	defer ts1.Close()

	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/get") {
			ts2SetHandler(w, r)
		} else if strings.HasPrefix(r.RequestURI, "/set") {
			ts2GetHandler(w, r)
		}
	}))

	defer ts2.Close()
	//////
	//	addrs := map[int]string{
	//		0: strings.TrimPrefix(ts1.URL, "http://"),
	//		1: strings.TrimPrefix(ts2.URL, "http://"),
	//	}
	//......
}
