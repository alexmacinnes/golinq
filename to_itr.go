package golinq

func IteratorToSlice[T_Out any](src Iterator[T_Out]) []T_Out {
	itr := src.initItr()

	result := []T_Out{}
	for {
		next, ok := itr.Next()
		if !ok {
			break
		}
		result = append(result, next)
	}

	return result
}

func IteratorToMap[T_In any, T_OutKey comparable, T_OutValue any](src Iterator[T_In], keyFunc func(T_In) T_OutKey, valueFunc func(T_In) T_OutValue) (map[T_OutKey]T_OutValue, bool) {
	itr := src.initItr()

	result := map[T_OutKey]T_OutValue{}

	for {
		next, ok := itr.Next()
		if !ok {
			break
		}

		key := keyFunc(next)
		_, exists := result[key]

		if exists {
			return nil, false
		}

		result[key] = valueFunc(next)
	}

	return result, true
}
