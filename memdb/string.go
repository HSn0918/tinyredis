package memdb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hsn/tiny-redis/RESP"
	"github.com/hsn/tiny-redis/logger"
)

func RegisterStringCommands() {
	RegisterCommand("set", setString)
	RegisterCommand("get", getString)
}
func setString(m *MemDb, cmd [][]byte) RESP.RedisData {
	cmdName := string(cmd[0])
	if strings.ToLower(cmdName) != "set" {
		logger.Error("setString Function: is not set")
		return RESP.MakeErrorData("server error ")

	}
	if len(cmd) < 3 {
		return RESP.MakeErrorData("error: commands is invalid")
	}
	// if key is expired,delete key
	m.CheckTTL(string(cmd[1]))
	var err error
	// nx 表示是否使用 SET 命令的 nx 选项，只在键不存在时设置键值对
	// xx 表示是否使用 SET 命令的 xx 选项，只在键已经存在时设置键值对
	// get 表示是否使用 SET 命令的 get 选项，如果为 true，需要返回原始键的值
	// ex 表示是否使用 SET 命令的 ex 选项，如果为 true，需要设置键的过期时间
	// keepttl 表示是否使用 SET 命令的 keepttl 选项，如果为 true，需要保持键的原有 TTL 不变
	var nx, xx, get, ex, keepttl bool
	var exval int64
	for i := 3; i < len(cmd); i++ {
		switch strings.ToLower(string(cmd[i])) {
		case "nx":
			nx = true
		case "xx":
			xx = true
		case "get":
			get = true
		case "keepttl":
			keepttl = true
		case "ex":
			ex = true
			i++
			if i >= len(cmd) {
				return RESP.MakeErrorData("error: commands is invalid")
			}
			exval, err = strconv.ParseInt(string(cmd[i]), 10, 64)
			if err != nil {
				return RESP.MakeErrorData(fmt.Sprintf("error: commands is invalid, %s is not interger", string(cmd[i])))
			}
		default:
			return RESP.MakeErrorData("Error unsupported option: " + string(cmd[i]))
		}
	}
	if (nx && xx) || (ex && keepttl) {
		return RESP.MakeErrorData("error: commands is invalid")
	}
	m.locks.Lock(string(cmd[1]))
	defer m.locks.UnLock(string(cmd[1]))
	var res RESP.RedisData

	oldVal, oldOk := m.db.Get(string(cmd[1]))
	// check is string
	var oldTypeVal []byte
	var typeOk bool
	if oldOk {
		oldTypeVal, typeOk = oldVal.([]byte)
		if !typeOk {
			return RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
		}
	}
	if nx || xx {
		if nx {
			if oldOk {
				m.db.Set(string(cmd[1]), cmd[2])
				res = RESP.MakeStringData("OK")
			} else {
				res = RESP.MakeNullBulkData()
			}
		} else {
			if oldOk {
				m.db.Set(string(cmd[1]), cmd[2])
				res = RESP.MakeStringData("OK")
			} else {
				res = RESP.MakeNullBulkData()
			}
		}
	} else {
		m.db.Set(string(cmd[1]), cmd[2])
		res = RESP.MakeStringData("OK")
	}
	if get {
		if !oldOk {
			res = RESP.MakeNullBulkData()
		} else {
			res = RESP.MakeBulkData(oldTypeVal)
		}
	}
	if !keepttl {
		m.DelTTL(string(cmd[1]))
	}
	if ex {
		m.SetTTL(string(cmd[1]), exval+time.Now().Unix())
	}
	return res
}
func getString(m *MemDb, cmd [][]byte) RESP.RedisData {
	if strings.ToLower(string(cmd[0])) != "get" {
		logger.Error("getString Function: cmdName is not get")
		return RESP.MakeErrorData("server error")
	}
	if len(cmd) != 2 {
		return RESP.MakeErrorData("error: commands is invalid")
	}
	key := string(cmd[1])
	if !m.CheckTTL(key) {
		return RESP.MakeNullBulkData()
	}
	m.locks.RLock(key)
	defer m.locks.RUnLock(key)
	val, ok := m.db.Get(key)
	if !ok {
		return RESP.MakeNullBulkData()
	}
	byteVal, ok := val.([]byte)
	if !ok {
		return RESP.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	return RESP.MakeBulkData(byteVal)
}
