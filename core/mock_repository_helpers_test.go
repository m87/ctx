package core

func copyPointerSlice[T any](items []*T) []*T {
	result := make([]*T, len(items))
	copy(result, items)
	return result
}

func limited[T any](items []T, limit int) []T {
	if limit > 0 && len(items) > limit {
		return items[:limit]
	}
	return items
}
