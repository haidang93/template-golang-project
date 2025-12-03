package arrayutil

func Where[T any](data *[]T, check func(*T) bool) []T {
	var res []T

	for _, n := range *data {
		if check(&n) {
			res = append(res, n)
		}
	}

	return res
}

func Find[T any](data *[]T, check func(*T) bool) *T {

	for _, n := range *data {
		if check(&n) {
			return &n
		}
	}

	return nil
}

func Map[T any, K any](data *[]K, convert func(*K) T) []T {
	var res []T

	for _, n := range *data {
		res = append(res, convert(&n))
	}

	return res
}

func ContainsPtr[T any](data *[]T, value *T) bool {
	for i := range *data {
		if &(*data)[i] == value {
			return true
		}
	}
	return false
}

func Contains[T comparable](data *[]T, value T) bool {
	if data == nil {
		return false
	}
	for _, v := range *data {
		if v == value {
			return true
		}
	}
	return false
}

func ArrayStringToArrayStringPtr(value *[]string) []*string {
	res := []*string{}
	for _, v := range *value {
		res = append(res, &v)
	}
	return res
}
