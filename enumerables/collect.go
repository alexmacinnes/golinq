package enumerables

func Any[T any](src Enumerable[T]) bool {
	resultChannel, cancelFunc := runAction(src)

	for x := range resultChannel {
		_ = x
		cancelFunc()
		return true
	}

	return false
}

func All[T any](src Enumerable[T], predicate func(T) bool) bool {
	resultChannel, cancelFunc := runAction(src)

	for x := range resultChannel {
		if !predicate(x) {
			cancelFunc()
			return false
		}
	}

	return true
}

func Contains[T comparable](src Enumerable[T], item T) bool {
	resultChannel, cancelFunc := runAction(src)

	for x := range resultChannel {
		if x == item {
			cancelFunc()
			return true
		}
	}

	return false
}

func ElementAt[T any](src Enumerable[T], index int) (T, bool) {
	resultChannel, cancelFunc := runAction(src)

	var result T
	var ok bool

	for i := 0; i < index; i++ {
		_, ok = consumeFirst(resultChannel)
		if !ok {
			return result, ok
		}
	}
	result, ok = consumeFirst(resultChannel)

	cancelFunc()

	return result, ok
}

func First[T any](src Enumerable[T]) (T, bool) {
	resultChannel, cancelFunc := runAction(src)

	result, ok := consumeFirst(resultChannel)
	cancelFunc()

	return result, ok
}

func FirstOrDefault[T any](src Enumerable[T]) T {
	resultChannel, cancelFunc := runAction(src)

	result, _ := consumeFirst(resultChannel)
	cancelFunc()

	return result
}

func Single[T any](src Enumerable[T]) (T, bool) {
	resultChannel, cancelFunc := runAction(src)

	result, ok := consumeFirst(resultChannel)

	if ok {
		// fail result if there is a second item in result channel
		_, nextOk := consumeFirst(resultChannel)
		if nextOk {
			var defaultValue T
			return defaultValue, false
		}
	}

	cancelFunc()

	return result, ok
}

func SingleOrDefault[T any](src Enumerable[T]) (T, bool) {
	resultChannel, cancelFunc := runAction(src)

	result, ok := consumeFirst(resultChannel)

	if ok {
		// fail result if there is a second item in result channel
		_, nextOk := consumeFirst(resultChannel)
		if nextOk {
			var defaultValue T
			return defaultValue, false
		}
	}

	cancelFunc()

	// to reach here there were 0 or 1 items, ok is true
	return result, true
}

func Last[T any](src Enumerable[T]) (T, bool) {
	resultChannel, _ := runAction(src)

	result, ok := consumeFirst(resultChannel)

	if !ok {
		return result, ok
	}

	for x := range resultChannel {
		result = x
	}

	return result, true
}

func LastOrDefault[T any](src Enumerable[T]) T {
	resultChannel, _ := runAction(src)

	var result T

	for x := range resultChannel {
		result = x
	}

	return result
}
