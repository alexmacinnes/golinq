package golinq

func consumeFirst[T any](src chan T) (T, bool) {
	for x := range src {
		return x, true
	}

	var defaultVal T
	return defaultVal, false
}

func Max[T Ordered](src Enumerable[T]) (T, bool) {
	resultChannel, _ := runAction(src)

	max, ok := consumeFirst(resultChannel)
	if !ok {
		return max, ok
	}

	for x := range resultChannel {
		if max < x {
			max = x
		}
	}

	return max, true
}

func Min[T Ordered](src Enumerable[T]) (T, bool) {
	resultChannel, _ := runAction(src)

	min, ok := consumeFirst(resultChannel)
	if !ok {
		return min, ok
	}

	for x := range resultChannel {
		if min > x {
			min = x
		}
	}

	return min, true
}

func Avg[T Numeric](src Enumerable[T]) (float64, bool) {
	resultChannel, _ := runAction(src)

	var total float64 = 0
	count := 0

	for x := range resultChannel {
		total += float64(x)
		count++
	}

	if count == 0 {
		return total, false
	}
	return total / float64(count), true
}

func Sum[T Numeric](src Enumerable[T]) T {
	resultChannel, _ := runAction(src)

	var total T = 0

	for x := range resultChannel {
		total += x
	}

	return total
}

func Count[T any](src Enumerable[T]) uint32 {
	resultChannel, _ := runAction(src)

	count := uint32(0)

	for x := range resultChannel {
		_ = x
		count++
	}

	return count
}

func Accumulate[TAccumulate any, TItem any](src Enumerable[TItem], seed TAccumulate, accumulator func(TAccumulate, TItem) TAccumulate) TAccumulate {
	resultChannel, _ := runAction(src)

	result := seed

	for x := range resultChannel {
		_ = x
		result = accumulator(result, x)
	}

	return result
}
