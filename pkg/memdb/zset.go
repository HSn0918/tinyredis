package memdb

import (
	RESP2 "github.com/hsn/tiny-redis/pkg/RESP"
	"github.com/hsn/tiny-redis/pkg/logger"
	"strconv"
	"strings"
)

func RegisterZSetCommands() {
	RegisterCommand("zadd", zAddZset)

}

func zAddZset(m *MemDb, cmd [][]byte) RESP2.RedisData {
	if strings.ToLower(string(cmd[0])) != "zadd" {
		logger.Error("zAddZset Function: cmdName is not zadd")
		return nil
	}
	if len(cmd) < 3 && len(cmd)%2 != 0 {
		return RESP2.MakeErrorData("wrong number of arguments for 'zadd' command")
	}
	key := string(cmd[1])
	m.CheckTTL(key)
	m.locks.Lock(key)
	defer m.locks.UnLock(key)
	temp, ok := m.db.Get(key)
	if !ok {
		temp = NewZSet()
		m.db.Set(key, temp)
	}
	zset, ok := temp.(*ZSet)
	if !ok {
		return RESP2.MakeErrorData("WRONGTYPE Operation against a key holding the wrong kind of value")
	}
	res := 0
	for i := 2; i < len(cmd); i += 2 {
		member := string(cmd[i+1])
		score, err := strconv.ParseFloat(string(cmd[i]), 64)
		if err != nil {
			return RESP2.MakeErrorData("ERR value is not a valid float")
		}
		zset.Add(member, score)
		res++
	}
	return RESP2.MakeIntData(int64(res))
}
