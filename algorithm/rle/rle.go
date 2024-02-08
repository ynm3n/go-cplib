package rle

func RunLengthEncoding[T comparable](s []T) (res []struct {
	value T
	count int
}) {
	for _, v := range s {
		if n := len(res); (n > 0) && res[n-1].value == v {
			res[n-1].count++
		} else {
			res = append(res, struct {
				value T
				count int
			}{v, 1})
		}
	}
	return res
}
