package mediacontentservice

// data utilities
func Filter[T any](items []T, condition func(item T) bool) []T {
	var new_items = make([]T, 0, len(items))
	for _, item := range items {
		if condition(item) {
			new_items = append(new_items, item)
		}
	}
	return new_items
}

func ForEach[T any](items []T, do func(item *T)) []T {
	for i := range items {
		do(&items[i])
	}
	return items
}

func Extract[T_input, T_output any](items []T_input, convert func(item *T_input) T_output) []T_output {
	var new_items = make([]T_output, len(items))
	for i := range items {
		new_items[i] = convert(&items[i])
	}
	return new_items
}

func Reduce[T any](items []T, reduce func(a, b T) T) T {
	var res T
	for i := range items {
		res = reduce(res, items[i])
	}
	return res
}

func In[T any](item T, list []T, compare func(a, b *T) bool) bool {
	return Index[T](item, list, compare) >= 0
}

func Any[T any](list []T, condition func(item *T) bool) bool {
	return IndexAny[T](list, condition) >= 0
}

func Index[T any](item T, list []T, compare func(a, b *T) bool) int {
	for i := range list {
		if compare(&item, &list[i]) {
			return i
		}
	}
	return -1
}

func IndexAny[T any](list []T, condition func(item *T) bool) int {
	for i := range list {
		if condition(&list[i]) {
			return i
		}
	}
	return -1
}

func SafeSlice[T any](array []T, start, noninclusive_end int) []T {
	if start < 0 {
		start = 0
	}
	if noninclusive_end < 0 {
		noninclusive_end = 0
	}
	return array[min(start, len(array)):min(noninclusive_end, len(array))]
}

// string utilities
func truncateTextWithEllipsis(text string, max_len int) string {
	if len(text) > max_len {
		return text[:max_len] + "..."
	}
	return text
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
