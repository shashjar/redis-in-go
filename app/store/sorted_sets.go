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

func (sortedSet *SortedSet) removeMembers(members []string) int {
	numRemovedMembers := 0
	for _, member := range members {
		currentScore, ok := sortedSet.Scores[member]
		if !ok {
			continue
		}

		sortedSet.SkipList.Delete(currentScore, member)
		delete(sortedSet.Scores, member)
		numRemovedMembers += 1
	}
	return numRemovedMembers
}

// Returns the rank and score of a member in the sorted set, and a boolean
// indicating whether the member exists
func (sortedSet *SortedSet) getRankAndScore(member string) (int, float64, bool) {
	if _, ok := sortedSet.Scores[member]; !ok {
		return 0, 0, false
	}

	return sortedSet.SkipList.GetRank(member), sortedSet.Scores[member], true
}
