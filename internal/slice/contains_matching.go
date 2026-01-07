package slice

func ContainsMatching[T any](s []T, condition func(item T) bool) bool {
	return len(Filter(s, condition)) > 0
}