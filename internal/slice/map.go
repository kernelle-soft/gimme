package slice

func Map[T any, R any](items []T, fn func(T) R) []R {
    result := make([]R, len(items))
    for i, item := range items {
        result[i] = fn(item)
    }
    return result
}

func MapErr[T any, R any](items []T, fn func(T) (R, error)) ([]R, error) {
    result := make([]R, len(items))
    for i, item := range items {
        val, err := fn(item)
        if err != nil {
            return nil, err
        }
        result[i] = val
    }
    return result, nil
}