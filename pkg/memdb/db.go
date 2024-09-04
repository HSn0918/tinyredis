package memdb

import (
	RESP2 "github.com/hsn/tiny-redis/pkg/RESP"
	"github.com/hsn/tiny-redis/pkg/config"
	"github.com/hsn/tiny-redis/pkg/logger"
	"strings"
	"time"
)

// MemDb is the memory cache database
// All key:value pairs are stored in db
// All ttl keys are stored in ttlKeys
// locks is used to lock a key for db to ensure some atomic operations
type MemDb struct {
	db      *ConcurrentMap
	ttlKeys *ConcurrentMap
	locks   *Locks
}

func NewMemDb() *MemDb {
	return &MemDb{
		db:      NewConcurrentMap(config.Configures.ShardNum),
		ttlKeys: NewConcurrentMap(config.Configures.ShardNum),
		locks:   NewLocks(config.Configures.ShardNum * 2),
	}
}
func (m *MemDb) ExecCommand(cmd [][]byte) RESP2.RedisData {
	if len(cmd) == 0 {
		return nil
	}
	var res RESP2.RedisData
	cmdName := strings.ToLower(string(cmd[0]))
	command, ok := CmdTable[cmdName]
	if !ok {
		res = RESP2.MakeErrorData("error: unsupported command")
	} else {
		execFunc := command.executor
		res = execFunc(m, cmd)
	}
	return res
}

// CheckTTL checks ttl keys and delete expired keys
// return false if key is expired,else true
// Attention: Don't lock this function because it has called locks.Lock(key) for atomic deleting expired key.
// Otherwise, it will cause a deadlock.
func (m *MemDb) CheckTTL(key string) bool {
	ttl, ok := m.ttlKeys.Get(key)
	if !ok {
		return true
	}
	ttlTime := ttl.(int64)
	now := time.Now().Unix()
	if ttlTime > now {
		return true
	}

	m.locks.Lock(key)
	defer m.locks.UnLock(key)
	m.db.Delete(key)
	m.ttlKeys.Delete(key)
	return false
}

// SetTTL sets ttl for keys
// return bool to check if ttl set success
// return int to check if the key is a new ttl key
func (m *MemDb) SetTTL(key string, value int64) int {
	if _, ok := m.db.Get(key); !ok {
		logger.Debug("SetTTL: key not exist")
		return 0
	}
	m.ttlKeys.Set(key, value)
	return 1
}
func (m *MemDb) DelTTL(key string) int {
	return m.ttlKeys.Delete(key)
}
