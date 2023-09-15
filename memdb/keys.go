package memdb

import (
	"fmt"
	"strings"
	"time"

	"github.com/hsn/tiny-redis/util"

	"github.com/hsn/tiny-redis/RESP"
	"github.com/hsn/tiny-redis/logger"
)

// RegisterKeyCommand
// Register command
func RegisterKeyCommand() {
	RegisterCommand("ping", pingKeys)
	RegisterCommand("del", delKey)
	RegisterCommand("exists", existsKey)
	RegisterCommand("keys", keysKey)
	RegisterCommand("expire", expireKey)
	RegisterCommand("persist", persistKey)
	RegisterCommand("ttl", ttlKey)
	RegisterCommand("type", typeKey)
	RegisterCommand("rename", renameKey)
}

// pingKeys
// if ping return pong
func pingKeys(m *MemDb, cmd [][]byte) RESP.RedisData {
	cmdName := string(cmd[0])
	if strings.ToLower(cmdName) != "ping" {
		logger.Error("pingKeys Function: cmdName is not ping")
		return RESP.MakeErrorData("server error")
	}
	if len(cmd) > 2 {
		return RESP.MakeErrorData("error: command args number is invalid")
	}
	if len(cmd) == 1 {
		return RESP.MakeStringData("PONG")
	}
	return RESP.MakeBulkData(cmd[1])
}
func delKey(m *MemDb, cmd [][]byte) RESP.RedisData {
	cmdName := string(cmd[0])
	if strings.ToLower(cmdName) != "del" {
		logger.Error("delKey Function: cmdName is not del")
		return RESP.MakeErrorData("Protocol error: command is not del")
	}
	if !m.CheckTTL(string(cmd[1])) {
		return RESP.MakeIntData(int64(0))
	}
	dKeyCount := 0
	for _, key := range cmd[1:] {
		m.locks.Lock(string(key))
		dKeyCount += m.db.Delete(string(key))
		m.ttlKeys.Delete(string(key))
		m.locks.UnLock(string(key))
	}
	return RESP.MakeIntData(int64(dKeyCount))
}
func existsKey(m *MemDb, cmd [][]byte) RESP.RedisData {
	cmdName := string(cmd[0])
	if strings.ToLower(cmdName) != "exist" || len(cmd) < 2 {
		logger.Error("existsKey Function: cmdName is not exists")
		return RESP.MakeErrorData("protocol error: command is not exists")
	}
	eKeyCount := 0
	for _, keyByte := range cmd[1:] {
		key := string(keyByte)
		if m.CheckTTL(key) {
			m.locks.RLock(key)
			if _, ok := m.db.Get(key); ok {
				eKeyCount++
			}
			m.locks.RUnLock(key)
		}
	}

	return RESP.MakeIntData(int64(eKeyCount))
}
func keysKey(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "keys" || len(cmd) != 2 {
		logger.Error("keysKey Function: cmdName is not keys or cmd length is not 2")
		return RESP.MakeErrorData(fmt.Sprintf("error: keys function get invalid command %s %s", string(cmd[0]), string(cmd[1])))
	}
	res := make([]RESP.RedisData, 0)
	allKeys := m.db.Keys()
	pattern := string(cmd[1])
	for _, key := range allKeys {
		if m.CheckTTL(key) {
			if util.PatternMatch(pattern, key) {
				res = append(res, RESP.MakeBulkData([]byte(key)))
			}
		}
	}
	return RESP.MakeArrayData(res)
}

// todo:
func expireKey(m *MemDb, cmd [][]byte) RESP.RedisData {

	return RESP.MakeNullBulkData()
}
func persistKey(m *MemDb, cmd [][]byte) RESP.RedisData {
	cmdName := string(cmd[0])
	if strings.ToLower(cmdName) != "persist" || len(cmd) != 2 {
		logger.Error("persistKey Function: cmdName is not persist or command args number is invalid")
		return RESP.MakeErrorData("error: cmdName is not persist or command args number is invalid")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP.MakeIntData(int64(0))
	}
	m.locks.Lock(key)
	defer m.locks.UnLock(key)
	res := m.DelTTL(key)
	return RESP.MakeIntData(int64(res))
}
func ttlKey(m *MemDb, cmd [][]byte) RESP.RedisData {
	cmdName := string(cmd[0])
	if strings.ToLower(cmdName) != "ttl" || len(cmd) != 2 {
		logger.Error("ttlKey Function: cmdName is not ttl or command args number is invalid")
		return RESP.MakeErrorData("error: cmdName is not ttl or command args number is invalid")
	}
	key := string(cmd[1])
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	if _, ok := m.db.Get(key); !ok {
		return RESP.MakeIntData(int64(-2))
	}
	ttl, ok := m.ttlKeys.Get(key)
	if !ok {
		return RESP.MakeIntData(int64(-1))
	}
	now := time.Now().Unix()
	return RESP.MakeIntData(ttl.(int64) - now)

}

// todo:type
func typeKey(m *MemDb, cmd [][]byte) RESP.RedisData {
	return RESP.MakeNullBulkData()

}
func renameKey(m *MemDb, cmd [][]byte) RESP.RedisData {
	cmdName := string(cmd[0])
	if strings.ToLower(cmdName) != "rename" || len(cmd) != 3 {
		logger.Error("renameKey Function: cmdName is not rename or command args number is not invalid")
		return RESP.MakeErrorData("error: cmdName is not rename or command args number is not invalid")
	}
	oldName, newName := string(cmd[1]), string(cmd[2])
	if !m.CheckTTL(oldName) {
		return RESP.MakeErrorData(fmt.Sprintf("error: %s not exist", oldName))
	}
	m.locks.RLockMulti([]string{oldName, newName})
	defer m.locks.RUnLockMulti([]string{oldName, newName})
	oldValue, ok := m.db.Get(oldName)
	if !ok {
		return RESP.MakeErrorData(fmt.Sprintf("error: %s not exist", oldName))
	}
	m.db.Delete(oldName)
	m.ttlKeys.Delete(oldName)
	m.db.Delete(newName)
	m.ttlKeys.Delete(newName)
	m.db.Set(newName, oldValue)
	return RESP.MakeStringData("OK")
}
