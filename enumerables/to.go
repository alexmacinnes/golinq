package enumerables

func EnumerableToSlice[T_Out any](src Enumerable[T_Out]) []T_Out {
	resultChannel, _ := runAction(src)

	result := []T_Out{}
	for x := range resultChannel {
		result = append(result, x)
	}

	return result
}

func EnumerableToMap[T_In any, T_OutKey comparable, T_OutValue any](src Enumerable[T_In], keyFunc func(T_In) T_OutKey, valueFunc func(T_In) T_OutValue) (map[T_OutKey]T_OutValue, bool) {
	resultChannel, cancelFunc := runAction(src)

	result := map[T_OutKey]T_OutValue{}

	for x := range resultChannel {
		key := keyFunc(x)
		_, exists := result[key]

		if exists {
			cancelFunc()
			return nil, false
		}

		result[key] = valueFunc(x)
	}

	return result, true
}
