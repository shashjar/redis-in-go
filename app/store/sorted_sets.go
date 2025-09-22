package store

// Represents a Redis sorted set, a collection of unique strings (members) ordered
// by an associated float64 score. Some use cases for sorted sets include:
//   - Leaderboards
//   - Rate limiters
type SortedSet struct {
	Scores   map[string]float64
	SkipList *SkipList
}

func createEmptySortedSet(setKey string, kvs *KeyValueStore) *SortedSet {
	sortedSet := &SortedSet{Scores: make(map[string]float64), SkipList: NewSkipList()}
	kvs.data[setKey] = KeyValue{Value: sortedSet, Type: "zset"}
	return sortedSet
}

func (sortedSet *SortedSet) addMembers(memberScores map[string]float64) int {
	numNewMembers := 0
	for member, score := range memberScores {
		currentScore, ok := sortedSet.Scores[member]
		if ok && currentScore == score {
			continue
		} else if !ok {
			numNewMembers += 1
		}

		sortedSet.SkipList.Insert(score, member)
		sortedSet.Scores[member] = score
	}
	return numNewMembers
}
