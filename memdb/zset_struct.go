package memdb

import "math/rand"

type zSetNode struct {
	member  string
	score   float64
	forward []*zSetNode
}
type ZSet struct {
	header *zSetNode
	level  int
}

const MaxLevel = 32

func NewZSetNode() *zSetNode {
	return &zSetNode{
		forward: make([]*zSetNode, MaxLevel),
	}
}
func NewZSet() *ZSet {
	return &ZSet{
		header: NewZSetNode(),
		level:  1,
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

	// Traverse the skip list and keep track of the update pointers
	for i := z.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && current.forward[i].score < score {
			current = current.forward[i]
		}
		update[i] = current
	}

	// Generate a random level for the new node
	newLevel := randomLevel()

	// If the new level is greater than the current level of the skip list,
	// update the update pointers and set the new level for the skip list
	if newLevel > z.level {
		for i := z.level; i < newLevel; i++ {
			update[i] = z.header
		}
		z.level = newLevel
	}

	// Create a new node with the given member and score
	newNode := NewZSetNode()
	newNode.member = member
	newNode.score = score

	// Insert the new node into the skip list
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}
}
