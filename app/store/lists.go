package store

// Represents a Redis list, a linked list of string values. Redis lists are frequently used to:
//   - Implement stacks and queues.
//   - Build queue management for background worker systems.
type List struct {
	Entries []string
}

func reverseSlice(s []string) []string {
	result := make([]string, len(s))
	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		result[i], result[j] = s[j], s[i]
	}
	return result
}

func createEmptyList(listKey string, kvs *KeyValueStore) *List {
	list := &List{Entries: []string{}}
	kvs.data[listKey] = KeyValue{Value: list, Type: "list"}
	return list
}

func (list *List) appendElements(elements []string) int {
	list.Entries = append(list.Entries, elements...)
	return len(list.Entries)
}

func (list *List) prependElements(elements []string) int {
	list.Entries = append(elements, list.Entries...)
	return len(list.Entries)
}

func (list *List) getElementsInRange(startIndex int, stopIndex int) []string {
	length := len(list.Entries)

	if stopIndex >= length {
		stopIndex = length - 1
	}
	if startIndex > stopIndex {
		return []string{}
	}

	return list.Entries[startIndex : stopIndex+1]
}

func (list *List) popLeftElements(popCount int) []string {
	length := len(list.Entries)
	if popCount >= length {
		poppedElements := list.Entries
		list.Entries = []string{}
		return poppedElements
	}

	poppedElements := list.Entries[:popCount]
	list.Entries = list.Entries[popCount:]
	return poppedElements
}
