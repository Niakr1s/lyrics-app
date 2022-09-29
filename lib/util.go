package lib

func Filter[T any](input []T, condition func(T) bool) []T {
	res := []T{}
	for _, v := range input {
		if condition(v) {
			res = append(res, v)
		}
	}
	return res
}
