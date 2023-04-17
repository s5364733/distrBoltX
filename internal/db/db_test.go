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

// createTempDb 创建temp_db文件
func createTempDb(t *testing.T, readOnly bool) *Database {
	t.Helper()

	f, err := os.CreateTemp(os.TempDir(), "kvdb")
	if err != nil {
		t.Fatalf("Could not create temp file: %v", err)
	}
	name := f.Name()
	f.Close()
	t.Cleanup(func() { os.Remove(name) })

	db, closeFunc, err := NewDatabase(name, readOnly)
	if err != nil {
		t.Fatalf("Could not create a new database: %v", err)
	}
	t.Cleanup(func() { closeFunc() })

	return db
}

func TestGetSet(t *testing.T) {
	db := createTempDb(t, false)

	if err := db.SetKey("party", []byte("Great")); err != nil {
		t.Fatalf("Could not write key: %v", err)
	}

	value, err := db.GetKey("party")
	if err != nil {
		t.Fatalf(`Could not get the key "party": %v`, err)
	}

	if !bytes.Equal(value, []byte("Great")) {
		t.Errorf(`Unexpected value for key "party": got %q, want %q`, value, "Great")
	}

	k, v, err := db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if !bytes.Equal(k, []byte("party")) || !bytes.Equal(v, []byte("Great")) {
		t.Errorf(`GetNextKeyForReplication(): got %q, %q; want %q, %q`, k, v, "party", "Great")
	}
}

func TestDeleteReplicationKey(t *testing.T) {
	db := createTempDb(t, false)

	setKey(t, db, "party", "Great")
	//setKey(t, db, "party1", "Great1")
	k, v, err := db.GetNextKeyForReplication()
	//k1, v1, err := db.GetNextKeyForReplication()
	//t.Logf("key,k1 %q  v1 %q  err: %v", k1, v1, err)
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if !bytes.Equal(k, []byte("party")) || !bytes.Equal(v, []byte("Great")) {
		t.Errorf(`GetNextKeyForReplication(): got %q, %q; want %q, %q`, k, v, "party", "Great")
	}

	if err := db.DeleteReplicationKey([]byte("party"), []byte("Bad")); err == nil {
		t.Fatalf(`DeleteReplicationKey("party", "Bad"): got nil error, want non-nil error`)
	}

	if err := db.DeleteReplicationKey([]byte("party"), []byte("Great")); err != nil {
		t.Fatalf(`DeleteReplicationKey("party", "Great"): got %q, want nil error`, err)
	}

	k, v, err = db.GetNextKeyForReplication()
	if err != nil {
		t.Fatalf(`Unexpected error for GetNextKeyForReplication(): %v`, err)
	}

	if k == nil || v == nil {
		t.Errorf(`GetNextKeyForReplication(): got %v, %v; want nil, nil`, k, v)
	}
}

func TestSetReadOnly(t *testing.T) {
	db := createTempDb(t, true)

	if err := db.SetKey("party", []byte("Bad")); err == nil {
		t.Fatalf("SetKey(%q, %q): got nil error, want non-nil error", "party", []byte("Bad"))
	}
}

func setKey(t *testing.T, d *Database, key, value string) {
	t.Helper()

	if err := d.SetKey(key, []byte(value)); err != nil {
		t.Fatalf("SetKey(%q, %q) failed: %v", key, value, err)
	}
}

func getKey(t *testing.T, d *Database, key string) string {
	t.Helper()
	value, err := d.GetKey(key)
	if err != nil {
		t.Fatalf("GetKey(%q) failed: %v", key, err)
	}

	return string(value)
}

func TestDeleteExtraKeys(t *testing.T) {
	db := createTempDb(t, false)

	setKey(t, db, "party", "Great")
	setKey(t, db, "us", "CapitalistPigs")
	//db.db.View(func(tx *bolt.Tx) error {
	//	b := tx.Bucket([]byte("default"))
	//	key, value := b.Cursor().First()
	//	fmt.Sprintf("value %s %s", key, value)
	//	return errors.New("s")
	//})
	if err := db.DeleteExtraKeys(func(name string) bool { return name == "us" }); err != nil {
		t.Fatalf("Could not delete extra keys: %v", err)
	}

	if value := getKey(t, db, "party"); value != "Great" {
		t.Errorf(`Unexpected value for key "party": got %q, want %q`, value, "Great")
	}

	if value := getKey(t, db, "us"); value != "" {
		t.Errorf(`Unexpected value for key "us": got %q, want %q`, value, "")
	}
}
