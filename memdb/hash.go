package memdb

import (
	"strings"

	"github.com/hsn/tiny-redis/RESP"
	"github.com/hsn/tiny-redis/logger"
)

func RegisterHashCommands() {
	RegisterCommand("hdel", hDelHash)
	RegisterCommand("hexists", hExistsHash)
	RegisterCommand("hget", hGetHash)
	RegisterCommand("hgetall", hGetAllHash)
	//RegisterCommand("hincrby", hIncrByHash)
	//RegisterCommand("hincrbyfloat", hIncrByFloatHash)
	//RegisterCommand("hkeys", hKeysHash)
	RegisterCommand("hlen", hLenHash)
	RegisterCommand("hmget", hMGetHash)
	RegisterCommand("hset", hSetHash)
	//RegisterCommand("hsetnx", hSetNxHash)
	//RegisterCommand("hvals", hValsHash)
	//RegisterCommand("hstrlen", hStrLenHash)
	//RegisterCommand("hrandfield", hRandFieldHash)
}

func hLenHash(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "hlen" {
		logger.Error("hLenHash Function: cmdName is not hlen")
		return RESP.MakeErrorData("Server error")
	}
	if len(cmd) != 2 {
		return RESP.MakeErrorData("wrong number of arguments for 'hlen' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP.MakeIntData(0)
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)

	tem, ok := m.db.Get(key)
	if !ok {
		return RESP.MakeIntData(0)
	}
	hash, ok := tem.(*Hash)
	if !ok {
		RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	res := hash.Len()
	return RESP.MakeIntData(int64(res))
}
func hSetHash(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "hset" {
		logger.Error("hMSetHash Function: cmdName is not hset")
		return RESP.MakeErrorData("Server error")
	}
	if len(cmd) < 4 || len(cmd)&1 == 1 {
		return RESP.MakeErrorData("wrong number of arguments for 'hset' command")
	}
	key := string(cmd[1])
	m.CheckTTL(key)

	m.locks.Lock(key)
	defer m.locks.UnLock(key)
	var hash *Hash
	tem, ok := m.db.Get(key)
	if !ok {
		hash = NewHash()
		m.db.Set(key, hash)
	} else {
		hash, ok = tem.(*Hash)
		if !ok {
			return RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}
	for i := 2; i < len(cmd); i += 2 {
		field := string(cmd[i])
		value := cmd[i+1]
		hash.Set(field, value)
	}
	return RESP.MakeStringData("OK")
}
func hGetHash(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "hget" {
		logger.Error("hGetHash Function: command name is not hget")
		return RESP.MakeErrorData("Server error")
	}
	if len(cmd) != 3 {
		return RESP.MakeErrorData("wrong number of arguments for 'hget' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP.MakeNullBulkData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP.MakeNullBulkData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	res := hash.Get(string(cmd[2]))
	if len(res) == 0 {
		return RESP.MakeBulkData(nil)
	}
	return RESP.MakeBulkData(res)
}
func hMGetHash(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "hmget" {
		logger.Error("hMGetHash Function: cmdName is not hmget")
		return RESP.MakeErrorData("Server error")
	}
	if len(cmd) < 3 {
		return RESP.MakeErrorData("wrong number of arguments for 'hmget' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP.MakeEmptyArrayData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP.MakeEmptyArrayData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	res := make([]RESP.RedisData, 0, len(cmd)-2)
	for i := 2; i < len(cmd); i++ {
		field := string(cmd[i])
		data := hash.Get(field)
		if len(data) == 0 {
			res = append(res, RESP.MakeBulkData(nil))
		} else {
			res = append(res, RESP.MakeBulkData(data))
		}
	}
	return RESP.MakeArrayData(res)
}
func hDelHash(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "hdel" {
		logger.Error("hDelHash Function: cmdName is not hdel")
		return RESP.MakeErrorData("Server error")
	}
	if len(cmd) < 3 {
		return RESP.MakeErrorData("wrong number of arguments for 'hdel' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP.MakeIntData(0)
	}
	m.locks.Lock(key)
	defer m.locks.UnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP.MakeIntData(0)
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	defer func() {
		if hash.IsEmpty() {
			m.db.Delete(key)
			m.DelTTL(key)
		}
	}()
	res := 0
	for i := 2; i < len(cmd); i++ {
		res += hash.Del(string(cmd[i]))
	}

	return RESP.MakeIntData(int64(res))
}
func hExistsHash(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "hexist" {
		logger.Error("hExistsHash Function: cmdName is not hexist")
		return RESP.MakeErrorData("Server error")
	}
	if len(cmd) != 3 {
		RESP.MakeErrorData("wrong number of arguments for 'hexists' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP.MakeIntData(0)
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		RESP.MakeIntData(0)
	}
	hash, ok := tem.(*Hash)
	if !ok {
		RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	if hash.Exist(string(cmd[2])) {
		return RESP.MakeIntData(1)
	}
	return RESP.MakeIntData(0)
}
func hGetAllHash(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "hget" {
		logger.Error("hGetAllHash Function: cmdName is not hgetall")
		return RESP.MakeErrorData("Server error")
	}
	if len(cmd) != 2 {
		return RESP.MakeErrorData("wrong number of arguments for 'hgetall' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP.MakeEmptyArrayData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP.MakeEmptyArrayData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	table := hash.Table()

	res := make([]RESP.RedisData, 0, len(table)*2)
	for k, v := range table {
		res = append(res, RESP.MakeBulkData([]byte(k)), RESP.MakeBulkData(v))
	}
	return RESP.MakeArrayData(res)
}
