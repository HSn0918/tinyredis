package memdb

import (
	"bytes"
	"testing"
	"time"

	"github.com/hsn/tiny-redis/config"
)

func init() {
	config.Configures = &config.Config{ShardNum: 100}
}
func TestDelkey(t *testing.T) {
	memdb := NewMemDb()
	memdb.db.Set("a", "a")
	memdb.db.Set("b", "b")
	memdb.ttlKeys.Set("b", time.Now().Unix()+10)

	del_a := delKey(memdb, [][]byte{[]byte("del"), []byte("a"), []byte("b")})

	if !bytes.Equal(del_a.ToBytes(), []byte(":2\r\n")) {
		t.Error("del reply is not correct")
	}

	_, ok1 := memdb.db.Get("a")
	_, ok2 := memdb.db.Get("b")
	_, ok3 := memdb.ttlKeys.Get("b")
	if ok1 || ok2 || ok3 {
		t.Error("del failed")
	}
}
