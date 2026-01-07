package slice

func From[T any](s ...[]T) []T {
    var total int
    for _, sliceItem := range s {
        total += len(sliceItem)
    }
    
    result := make([]T, 0, total)
    for _, sliceItem := range s {
        result = append(result, sliceItem...)
    }
    return result
}