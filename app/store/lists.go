package store

// Represents a Redis list, a linked list of string values. Redis lists are frequently used to:
//   - Implement stacks and queues.
//   - Build queue management for background worker systems.
type List struct {
	Entries []string
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

func (list *List) getElementsInRange(startIndex int, stopIndex int) []string {
	return list.Entries[startIndex : stopIndex+1]
}
