package memdb

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcurrentMap_ManyHashCollisions(t *testing.T) {
	// 使用较小的哈希表大小以增加冲突概率
	cmap := NewConcurrentMap(4) // 将 size 设置为较小值

	// 生成大量键值对并插入
	numKeys := 100
	keys := make([]string, numKeys)

	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("key_%d", i)
		cmap.Set(keys[i], i)
	}

	// 计算实际产生冲突的键
	collisionCount := 0
	shardCounts := make(map[int]int)

	for _, key := range keys {
		pos := cmap.getKeyPos(key)
		shardCounts[pos]++
		if shardCounts[pos] > 1 {
			collisionCount++
		}
	}

	// 期望有一些冲突
	assert.Greater(t, collisionCount, 0, "There should be some hash collisions")

	// 确保所有键都能正确获取
	for i, key := range keys {
		value, found := cmap.Get(key)
		assert.True(t, found, "Value should be found for key")
		assert.Equal(t, i, value, "Value should match the inserted value")
	}

	// 删除冲突键并验证删除操作
	for _, key := range keys {
		deleted := cmap.Delete(key)
		assert.Equal(t, 1, deleted, "Key should be successfully deleted")
		value, found := cmap.Get(key)
		assert.False(t, found, "Key should not be found after deletion")
		assert.Nil(t, value, "Deleted key should return nil value")
	}
}
func TestConcurrentMap_ManyHashCollisions_ListNode(t *testing.T) {
	// 使用较小的哈希表大小以增加冲突概率
	cmap := NewConcurrentMap(4) // 将 size 设置为较小值

	// 生成大量键值对并插入
	numKeys := 10
	keys := make([]string, numKeys)

	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("key_%d", i)
		cmap.Set(keys[i], i)
	}

	// 计算实际产生冲突的键
	collisionCount := 0
	shardCounts := make(map[int]int)

	for _, key := range keys {
		pos := cmap.getKeyPos(key)
		shardCounts[pos]++
		if shardCounts[pos] > 1 {
			collisionCount++
		}
	}

	// 期望有一些冲突
	assert.Greater(t, collisionCount, 0, "There should be some hash collisions")

	// 确保所有键都能正确获取
	for i, key := range keys {
		value, found := cmap.Get(key)
		assert.True(t, found, "Value should be found for key")
		assert.Equal(t, i, value, "Value should match the inserted value")
	}

	// 删除冲突键并验证删除操作
	for _, key := range keys {
		deleted := cmap.Delete(key)
		assert.Equal(t, 1, deleted, "Key should be successfully deleted")
		value, found := cmap.Get(key)
		assert.False(t, found, "Key should not be found after deletion")
		assert.Nil(t, value, "Deleted key should return nil value")
	}
}
