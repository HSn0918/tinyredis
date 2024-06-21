package memdb

import (
	"fmt"
	"math/rand"
)

type zSetNode struct {
	member  string
	score   float64
	forward []*zSetNode
}
type ZSet struct {
	header *zSetNode
	level  int
	length int
}

func (z *ZSet) String() string {
	var result string
	for i := z.level - 1; i >= 0; i-- {
		result += fmt.Sprintf("Level %d: ", i+1)
		current := z.header.forward[i]
		for current != nil {
			result += fmt.Sprintf("%s(%f) ", current.member, current.score)
			current = current.forward[i]
		}
		result += "\n"
	}
	return result
}

const MaxLevel = 32

func NewZSetNode(level int, member string, score float64) *zSetNode {
	return &zSetNode{
		member:  member,
		score:   score,
		forward: make([]*zSetNode, level),
	}
}
func NewZSet() *ZSet {
	return &ZSet{
		header: NewZSetNode(MaxLevel, "", 0),
		level:  1,
		length: 0, // 初始化长度为 0
	}
}
func randomLevel() int {
	level := 1
	for rand.Float64() < 0.25 && level < MaxLevel {
		level++
	}
	return level
}

func (z *ZSet) Add(member string, score float64) {
	update := make([]*zSetNode, MaxLevel)
	current := z.header
	for i := z.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].score < score {
			current = current.forward[i]
		}
		update[i] = current
	}
	newLevel := randomLevel()
	if newLevel > z.level {
		for i := z.level; i < newLevel; i++ {
			update[i] = z.header
		}
		z.level = newLevel
	}
	newNode := NewZSetNode(newLevel, member, score)
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
	z.length++
}
func (z *ZSet) Get(member string) float64 {
	current := z.header
	for i := z.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].member < member {
			current = current.forward[i]
		}
	}
	current = current.forward[0]
	if current != nil && current.member == member {
		return current.score
	}
	return -1
}
func (z *ZSet) Remove(member string, score float64) bool {
	update := make([]*zSetNode, MaxLevel)
	current := z.header

	for i := z.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].score < score {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	if current != nil && current.score == score && current.member == member {
		for i := 0; i < z.level; i++ {
			if update[i].forward[i] != current {
				break
			}
			update[i].forward[i] = current.forward[i]
		}

		for z.level > 1 && z.header.forward[z.level-1] == nil {
			z.level--
		}

		z.length--
		return true
	}

	return false
}

func (z *ZSet) Len() int {
	return z.length
}
