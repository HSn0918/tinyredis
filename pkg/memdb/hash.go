package memdb

import (
	"fmt"
	RESP2 "github.com/hsn/tiny-redis/pkg/RESP"
	"github.com/hsn/tiny-redis/pkg/logger"
	"strconv"
	"strings"
)

func RegisterHashCommands() {
	RegisterCommand("hdel", hDelHash)
	RegisterCommand("hexists", hExistsHash)
	RegisterCommand("hget", hGetHash)
	RegisterCommand("hgetall", hGetAllHash)
	RegisterCommand("hincrby", hIncrByHash)
	RegisterCommand("hincrbyfloat", hIncrByFloatHash)
	RegisterCommand("hkeys", hKeysHash)
	RegisterCommand("hlen", hLenHash)
	RegisterCommand("hmget", hMGetHash)
	RegisterCommand("hset", hSetHash)
	RegisterCommand("hsetnx", hSetNxHash)
	RegisterCommand("hvals", hValsHash)
	RegisterCommand("hstrlen", hStrLenHash)
	RegisterCommand("hrandfield", hRandFieldHash)
}

func hRandFieldHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hrandfield" {
		logger.Error("hRandFieldHash Function: cmdName is not hrandfield")
		return RESP2.MakeErrorData("server error")
	}

	if len(cmd) != 2 && len(cmd) != 3 && len(cmd) != 4 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hrandfield' command")
	}
	key := string(cmd[1])
	count := 1
	withValues := false
	var err error
	if len(cmd) >= 3 {
		count, err = strconv.Atoi(string(cmd[2]))
		if err != nil {
			return RESP2.MakeErrorData("err: count value must be integer")
		}
	}
	if len(cmd) == 4 {
		if strings.ToLower(string(cmd[3])) == "withvalues" {
			withValues = true
		} else {
			return RESP2.MakeErrorData(fmt.Sprintf("command option error, not support %s option", string(cmd[3])))
		}
	}
	if !m.CheckTTL(key) {
		return RESP2.MakeEmptyArrayData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeEmptyArrayData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	res := make([]RESP2.RedisData, 0)
	if withValues {
		fields := hash.RandomWithValue(count)
		for _, v := range fields {
			res = append(res, RESP2.MakeBulkData(v))
		}
	} else {
		fields := hash.Random(count)
		for _, v := range fields {
			res = append(res, RESP2.MakeBulkData([]byte(v)))
		}
	}
	return RESP2.MakeArrayData(res)
}

func hStrLenHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hstrlen" {
		logger.Error("hStrLenHash Function: cmdName is not hstrlen")
		return RESP2.MakeErrorData("server error")
	}
	if len(cmd) != 3 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hstrlen' command")
	}
	key := string(cmd[1])
	field := string(cmd[2])

	if !m.CheckTTL(key) {
		return RESP2.MakeIntData(0)
	}

	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeIntData(0)
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	res := hash.StrLen(field)

	return RESP2.MakeIntData(int64(res))
}

func hSetNxHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hsetnx" {
		logger.Error("hSetNxHash Function: cmdName is not hsetnx")
		return RESP2.MakeErrorData("server error")
	}
	if len(cmd) != 4 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hsetnx' command")
	}
	key := string(cmd[1])
	field := string(cmd[2])
	value := cmd[3]
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
			return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	if hash.Exist(field) {
		return RESP2.MakeIntData(0)
	}

	hash.Set(field, value)
	return RESP2.MakeIntData(1)

}

func hIncrByFloatHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hincrbyfloat" {
		logger.Error("hIncrByFloatHash Function: cmdName is not hincrbyfloat")
		return RESP2.MakeErrorData("server error")
	}

	if len(cmd) != 4 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hincrbyfloat' command")
	}

	var hash *Hash
	key, field := string(cmd[1]), string(cmd[2])
	incr, err := strconv.ParseFloat(string(cmd[3]), 64)
	if err != nil {
		return RESP2.MakeErrorData("incr value must be a float")
	}
	m.CheckTTL(key)

	m.locks.Lock(key)
	defer m.locks.UnLock(key)

	tem, ok := m.db.Get(key)
	if !ok {
		hash = NewHash()
		m.db.Set(key, hash)
	} else {
		hash, ok = tem.(*Hash)
		if !ok {
			return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}

	res, ok := hash.IncrByFloat(field, incr)
	if !ok {
		return RESP2.MakeErrorData("value is not a float")
	}

	return RESP2.MakeBulkData([]byte(strconv.FormatFloat(res, 'f', -1, 64)))
}

func hIncrByHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hincrby" {
		logger.Error("hIncrByHash Function: cmdName is not hincrby")
		return RESP2.MakeErrorData("server error")
	}
	if len(cmd) != 4 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hincrby' command")
	}
	var incr int
	var err error
	var hash *Hash
	key := string(cmd[1])
	field := string(cmd[2])
	incr, err = strconv.Atoi(string(cmd[3]))
	if err != nil {
		return RESP2.MakeErrorData("incr value must be an integer")
	}
	m.CheckTTL(key)
	m.locks.Lock(key)
	defer m.locks.UnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		hash = NewHash()
		m.db.Set(key, hash)
	} else {
		hash, ok = tem.(*Hash)
		if !ok {
			return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}
	res, ok := hash.IncrBy(field, incr)
	if !ok {
		return RESP2.MakeErrorData("value is not an integer")
	}
	return RESP2.MakeIntData(int64(res))
}

func hValsHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hvals" {
		logger.Error("hValsHash Function: cmdName is not hvals")
		return RESP2.MakeErrorData("server error")
	}
	if len(cmd) != 2 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hvals' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP2.MakeEmptyArrayData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeEmptyArrayData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	vals := hash.Values()
	res := make([]RESP2.RedisData, 0, len(vals))
	for _, val := range vals {
		res = append(res, RESP2.MakeBulkData(val))
	}
	return RESP2.MakeArrayData(res)
}

func hKeysHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hkeys" {
		logger.Error("hKeysHash Function: cmdName is not hkeys")
		return RESP2.MakeErrorData("server error")
	}
	if len(cmd) != 2 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hkeys' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP2.MakeEmptyArrayData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeEmptyArrayData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	fields := hash.Keys()
	res := make([]RESP2.RedisData, 0, len(fields))
	for _, v := range fields {
		res = append(res, RESP2.MakeBulkData([]byte(v)))
	}
	return RESP2.MakeArrayData(res)
}
func hLenHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hlen" {
		logger.Error("hLenHash Function: cmdName is not hlen")
		return RESP2.MakeErrorData("Server error")
	}
	if len(cmd) != 2 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hlen' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP2.MakeIntData(0)
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)

	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeIntData(0)
	}
	hash, ok := tem.(*Hash)
	if !ok {
		RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	res := hash.Len()
	return RESP2.MakeIntData(int64(res))
}
func hSetHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hset" {
		logger.Error("hMSetHash Function: cmdName is not hset")
		return RESP2.MakeErrorData("Server error")
	}
	if len(cmd) < 4 || len(cmd)&1 == 1 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hset' command")
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
			return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}
	for i := 2; i < len(cmd); i += 2 {
		field := string(cmd[i])
		value := cmd[i+1]
		hash.Set(field, value)
	}
	return RESP2.MakeStringData("OK")
}
func hGetHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hget" {
		logger.Error("hGetHash Function: command name is not hget")
		return RESP2.MakeErrorData("Server error")
	}
	if len(cmd) != 3 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hget' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP2.MakeNullBulkData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeNullBulkData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	res := hash.Get(string(cmd[2]))
	if len(res) == 0 {
		return RESP2.MakeBulkData(nil)
	}
	return RESP2.MakeBulkData(res)
}
func hMGetHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hmget" {
		logger.Error("hMGetHash Function: cmdName is not hmget")
		return RESP2.MakeErrorData("Server error")
	}
	if len(cmd) < 3 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hmget' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP2.MakeEmptyArrayData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeEmptyArrayData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	res := make([]RESP2.RedisData, 0, len(cmd)-2)
	for i := 2; i < len(cmd); i++ {
		field := string(cmd[i])
		data := hash.Get(field)
		if len(data) == 0 {
			res = append(res, RESP2.MakeBulkData(nil))
		} else {
			res = append(res, RESP2.MakeBulkData(data))
		}
	}
	return RESP2.MakeArrayData(res)
}
func hDelHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hdel" {
		logger.Error("hDelHash Function: cmdName is not hdel")
		return RESP2.MakeErrorData("Server error")
	}
	if len(cmd) < 3 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hdel' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP2.MakeIntData(0)
	}
	m.locks.Lock(key)
	defer m.locks.UnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeIntData(0)
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
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

	return RESP2.MakeIntData(int64(res))
}
func hExistsHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hexist" {
		logger.Error("hExistsHash Function: cmdName is not hexist")
		return RESP2.MakeErrorData("Server error")
	}
	if len(cmd) != 3 {
		RESP2.MakeErrorData("wrong number of arguments for 'hexists' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP2.MakeIntData(0)
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		RESP2.MakeIntData(0)
	}
	hash, ok := tem.(*Hash)
	if !ok {
		RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	if hash.Exist(string(cmd[2])) {
		return RESP2.MakeIntData(1)
	}
	return RESP2.MakeIntData(0)
}
func hGetAllHash(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "hget" {
		logger.Error("hGetAllHash Function: cmdName is not hgetall")
		return RESP2.MakeErrorData("Server error")
	}
	if len(cmd) != 2 {
		return RESP2.MakeErrorData("wrong number of arguments for 'hgetall' command")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP2.MakeEmptyArrayData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	tem, ok := m.db.Get(key)
	if !ok {
		return RESP2.MakeEmptyArrayData()
	}
	hash, ok := tem.(*Hash)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	table := hash.Table()

	res := make([]RESP2.RedisData, 0, len(table)*2)
	for k, v := range table {
		res = append(res, RESP2.MakeBulkData([]byte(k)), RESP2.MakeBulkData(v))
	}
	return RESP2.MakeArrayData(res)
}
