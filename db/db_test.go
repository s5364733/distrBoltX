package db

import (
	"bytes"
	"os"
	"testing"
)

//func createConfig(t *testing.T, content string) config.Config {
//	t.Helper()
//
//	f, err := os.CreateTemp(os.TempDir(), "test.db")
//	if err != nil {
//		t.Fatalf("Coundn't create the temp file : %v", err)
//	}
//	defer f.Close()
//	name := f.Name()
//	defer os.Remove(name)
//
//	_, err = f.WriteString(content)
//	if err != nil {
//		t.Fatalf("Write file occurs error ,:%v", err)
//	}
//
//	c, err := ParseConfig(name)
//	if err != nil {
//		t.Fatalf("Couldn't parse config :%v ", err)
//	}
//	return c
//}

func TestGetSet(t *testing.T) {

	f, err := os.CreateTemp(os.TempDir(), "kvdb")
	if err != nil {
		t.Fatalf("Could not create temp file :%v", err)
	}

	name := f.Name()
	f.Close()
	defer os.Remove(name)

	db, closeFunc, err := NewDatabase(name)
	if err != nil {
		t.Fatalf("Could create new db %v", err)
	}
	defer closeFunc()

	setKey(t, db, "key", "value")
	//
	//if err := db.SetKey("key", []byte("value")); err != nil {
	//	t.Fatalf("Could not create write key %v", err)
	//}

	value, err := db.GetKey("key")
	if err != nil {
		t.Fatalf("Could not get key %v", err)
	}

	if !bytes.Equal(value, []byte("value")) {
		t.Fatalf("Unexpected value for key 'key' value : %q", value)
	}
}

func setKey(t *testing.T, d *Database, key, value string) {
	t.Helper()
	if err := d.SetKey(key, []byte(value)); err != nil {
		t.Fatalf("Setkey(%q,%q)(failed) %v", key, value, err)
	}
}

func getKey(t *testing.T, d *Database, key string) string {
	t.Helper()
	val, err := d.GetKey(key)
	if err != nil {
		t.Fatalf("GetKey(%q)(failed) %v", key, err)
	}
	return string(val)
}

func TestDatabase_DeleteExtraKeys(t *testing.T) {

	f, err := os.CreateTemp(os.TempDir(), "kvdb")
	if err != nil {
		t.Fatalf("Could not create temp file :%v", err)
	}

	name := f.Name()
	f.Close()
	defer os.Remove(name)

	db, closeFunc, err := NewDatabase(name)
	if err != nil {
		t.Fatalf("Could create new db %v", err)
	}
	defer closeFunc()

	setKey(t, db, "key", "value")
	setKey(t, db, "us", "great")

	if err := db.DeleteExtraKeys(func(name string) bool { return name == "us" }); err != nil {
		t.Fatalf("Could delate extra keys %v", err)
	}

	if v := getKey(t, db, "key"); v != "value" {
		t.Fatalf("Unexpected value for key 'key' value : %q", v)
	}

	if v := getKey(t, db, "us"); v != "" {
		t.Fatalf("Unexpected value for key 'key' value : %q", v)
	}

}
