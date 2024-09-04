package memdb

import (
	"fmt"
	"testing"
)

func TestZSet_Add(t *testing.T) {
	zset := NewZSet()

	for i := range 128 {
		zset.Add(fmt.Sprintf("number%d", i), float64(i))
	}
	if zset.Len() != 128 {
		t.Errorf("expected length 3, got %d", zset.Len())
	}
	t.Logf("zset: %+v", zset)
}

func TestZSet_Get(t *testing.T) {
	zset := NewZSet()

	zset.Add("member1", 1.0)
	value := zset.Get("member1")
	if value != 1.0 {
		t.Errorf("expected value 1.0, got %f", value)
	}

}

func TestZSet_Len(t *testing.T) {
	zset := NewZSet()

	if zset.Len() != 0 {
		t.Errorf("expected length 0, got %d", zset.Len())
	}

	zset.Add("member1", 1.0)
	if zset.Len() != 1 {
		t.Errorf("expected length 1, got %d", zset.Len())
	}
}
