package config

import (
	"os"
	"reflect"
	"testing"
)

func createConfig(t *testing.T, content string) Config {
	t.Helper()

	f, err := os.CreateTemp(os.TempDir(), "config.toml")
	if err != nil {
		t.Fatalf("Coundn't create the temp file : %v", err)
	}
	defer f.Close()
	name := f.Name()
	defer os.Remove(name)

	_, err = f.WriteString(content)
	if err != nil {
		t.Fatalf("Write file occurs error ,:%v", err)
	}

	c, err := ParseConfig(name)
	if err != nil {
		t.Fatalf("Couldn't parse config :%v ", err)
	}
	return c
}

func TestParseConfig(t *testing.T) {
	got := createConfig(t, `[[shards]]
    name = "shard0"
    idx  = 0
    address = "localhost:8080"
   `)

	expect := &Config{
		Shards: []Shard{
			{
				Name:    "shard0",
				Idx:     0,
				Address: "localhost:8080",
			},
		},
	}
	//I really fuck expect a pointer type for an hour
	//!reflect.DeepEqual(c, expect)
	if !reflect.DeepEqual(got, *expect) {
		t.Errorf("The config not match source:%#v expect :%#v", got, expect)
	}
}

func TestParseShards(t *testing.T) {
	c := createConfig(t, `
[[shards]]
    name = "shard0"
    idx  = 0
    address = "localhost:8080"
[[shards]]
	name = "shard1"
	idx  = 1
	address = "localhost:8081"
   
`)
	got, err := ParseShards(c.Shards, "shard1")

	if err != nil {
		t.Fatalf("Cound not parse shards %#v:%v", c.Shards, err)
	}

	expect := &Shards{
		Count:  2,
		CurIdx: 1,
		Addrs: map[int]string{
			0: "localhost:8080",
			1: "localhost:8081",
		},
	}
	if !reflect.DeepEqual(got, expect) {
		t.Errorf("The shards config does match source:%#v expect :%#v", got, expect)
	}
}
