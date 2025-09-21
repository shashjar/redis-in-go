package store

import "math/rand/v2"

// Represents a skip list data structure
type SkipList struct {
	head     *SkipListNode
	maxLevel int
	level    int // Current highest level
	size     int // Number of elements
}

// Represents a node in the skip list
type SkipListNode struct {
	member string
	score  float64
	next   []*SkipListNode // Array of pointers to next nodes at each level
}

const (
	// Maximum level for the skip list (can be adjusted based on expected size)
	MAX_LEVEL = 16
	// Probability for level promotion (0.5 = 50%)
	P = 0.5
)

// Creates a new skip list
func NewSkipList() *SkipList {
	head := &SkipListNode{
		member: "",
		score:  0,
		next:   make([]*SkipListNode, MAX_LEVEL+1),
	}

	return &SkipList{
		head:     head,
		maxLevel: MAX_LEVEL,
		level:    0,
		size:     0,
	}
}

// Finds a node with the given member
// TODO: currently doing linear search, improve this
func (sl *SkipList) Search(member string) *SkipListNode {
	current := sl.head.next[0]
	for current != nil {
		if current.member == member {
			return current
		}
		current = current.next[0]
	}
	return nil
}

// Adds a new node with the given score and member to the skip list
func (sl *SkipList) Insert(score float64, member string) *SkipListNode {
	existingNode := sl.Search(member)
	if existingNode != nil {
		// If score is the same, no need to reinsert
		if existingNode.score == score {
			return existingNode
		}

		// Score changed, need to delete and reinsert at correct position
		sl.Delete(existingNode.score, member)
	}

	// Find insertion point and track previous nodes at each level
	update := make([]*SkipListNode, sl.maxLevel+1)
	current := sl.head

	// Find the position to insert
	for i := sl.level; i >= 0; i-- {
		for current.next[i] != nil &&
			(current.next[i].score < score ||
				(current.next[i].score == score && current.next[i].member < member)) {
			current = current.next[i]
		}
		update[i] = current
	}

	// Generate random level for new node
	newLevel := sl.randomLevel()

	// If new level is higher than current level, update head pointers
	if newLevel > sl.level {
		for i := sl.level + 1; i <= newLevel; i++ {
			update[i] = sl.head
		}
		sl.level = newLevel
	}

	// Create new node
	newNode := &SkipListNode{
		member: member,
		score:  score,
		next:   make([]*SkipListNode, newLevel+1),
	}

	// Update pointers
	for i := 0; i <= newLevel; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode
	}

	sl.size++
	return newNode
}

// Removes a node with the given score and member from the skip list
func (sl *SkipList) Delete(score float64, member string) bool {
	// Find the node to delete and track previous nodes at each level
	update := make([]*SkipListNode, sl.maxLevel+1)
	current := sl.head

	for i := sl.level; i >= 0; i-- {
		for current.next[i] != nil &&
			(current.next[i].score < score ||
				(current.next[i].score == score && current.next[i].member < member)) {
			current = current.next[i]
		}
		update[i] = current
	}

	// Check if node exists
	current = current.next[0]
	if current == nil || current.score != score || current.member != member {
		return false
	}

	// Update pointers to skip the deleted node
	for i := 0; i <= sl.level; i++ {
		if update[i].next[i] != current {
			break
		}
		update[i].next[i] = current.next[i]
	}

	// Update level if necessary
	for sl.level > 0 && sl.head.next[sl.level] == nil {
		sl.level--
	}

	sl.size--
	return true
}

func (sl *SkipList) randomLevel() int {
	level := 0
	for rand.Float64() < P && level < sl.maxLevel {
		level++
	}
	return level
}
