package memdb

import (
	"github.com/hsn/tiny-redis/pkg/logger"
	"github.com/hsn/tiny-redis/pkg/util"
	"sort"
	"sync"
)

type Locks struct {
	locks []*sync.RWMutex
}

func NewLocks(size int) *Locks {
	locks := make([]*sync.RWMutex, size)
	for i := 0; i < size; i++ {
		locks[i] = &sync.RWMutex{}
	}
	return &Locks{locks: locks}
}
func (l *Locks) GetKeyPos(key string) int {
	if len(l.locks) == 0 {
		// 防止除零错误
		logger.Error("Locks array is empty, cannot get key position")
		return -1
	}

	pos := util.HashKey(key)
	if pos < 0 {
		pos = -pos // 确保 pos 为正数
	}

	return pos % len(l.locks)
}
func (l *Locks) Lock(key string) {
	pos := l.GetKeyPos(key)
	if pos == -1 {
		logger.Error("Locks Lock key %s error: pos == -1", key)
		return
	}
	l.locks[pos].Lock()
}
func (l *Locks) UnLock(key string) {
	pos := l.GetKeyPos(key)
	if pos == -1 {
		logger.Error("Locks UnLock key %s error: pos == -1", key)
	}
	l.locks[pos].Unlock()
}
func (l *Locks) RLock(key string) {
	pos := l.GetKeyPos(key)
	if pos == -1 {
		logger.Error("Locks RLock key %s error: pos == -1", key)
	}
	l.locks[pos].RLock()
}
func (l *Locks) RUnLock(key string) {
	pos := l.GetKeyPos(key)
	if pos == -1 {
		logger.Error("Locks RUnLock key %s error: pos == -1", key)
	}
	l.locks[pos].RUnlock()
}
func (l *Locks) sortedLockPoses(keys []string) []int {
	set := make(map[int]struct{})
	for _, key := range keys {
		pos := l.GetKeyPos(key)
		if pos == -1 {
			logger.Error("Locks Lock key %s error: pos == -1", key)
			return nil
		}
		set[pos] = struct{}{}
	}
	poses := make([]int, len(set))
	i := 0
	for pos := range poses {
		poses[i] = pos
		i++
	}
	sort.Ints(poses)
	return poses
}
func (l *Locks) LockMulti(keys []string) {
	poses := l.sortedLockPoses(keys)
	if poses == nil {
		return
	}
	for _, pos := range poses {
		l.locks[pos].Lock()
	}
}
func (l *Locks) UnLockMulti(keys []string) {
	poses := l.sortedLockPoses(keys)
	if poses == nil {
		return
	}
	for _, pos := range poses {
		l.locks[pos].Unlock()
	}
}
func (l *Locks) RLockMulti(keys []string) {
	poses := l.sortedLockPoses(keys)
	if poses == nil {
		return
	}
	for _, pos := range poses {
		l.locks[pos].RLock()
	}
}
func (l *Locks) RUnLockMulti(keys []string) {
	poses := l.sortedLockPoses(keys)
	if poses == nil {
		return
	}
	for _, pos := range poses {
		l.locks[pos].RUnlock()
	}
}
