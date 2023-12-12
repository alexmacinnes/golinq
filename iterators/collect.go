package iterators

func AnyItr[T any](src Iterator[T]) bool {
	itr := src.initItr()
	_, ok := itr.Next()
	return ok
}

func AllItr[T any](src Iterator[T], predicate func(T) bool) bool {
	itr := src.initItr()

	for {
		next, ok := itr.Next()
		if !ok {
			return true
		}
		match := predicate(next)
		if !match {
			return false
		}
	}
}

func ContainsItr[T comparable](src Iterator[T], item T) bool {
	itr := src.initItr()

	for {
		next, ok := itr.Next()
		if !ok {
			return false
		}
		if next == item {
			return true
		}
	}
}

func ElementAtItr[T any](src Iterator[T], index int) (T, bool) {
	itr := src.initItr()
	var result T
	var ok bool

	for i := 0; i < index; i++ {
		_, ok = itr.Next()
		if !ok {
			return result, ok
		}
	}
	result, ok = itr.Next()

	return result, ok
}

func FirstItr[T any](src Iterator[T]) (T, bool) {
	itr := src.initItr()
	return itr.Next()
}

func FirstOrDefaultItr[T any](src Iterator[T]) T {
	itr := src.initItr()
	result, _ := itr.Next()
	return result
}

func SingleItr[T any](src Iterator[T]) (T, bool) {
	itr := src.initItr()

	result, ok := itr.Next()

	if ok {
		// fail result if there is a second item in result channel
		_, nextOk := itr.Next()
		if nextOk {
			var defaultValue T
			return defaultValue, false
		}
	}

	return result, ok
}

func SingleOrDefaultItr[T any](src Iterator[T]) (T, bool) {
	itr := src.initItr()

	result, ok := itr.Next()

	if ok {
		// fail result if there is a second item in result channel
		_, nextOk := itr.Next()
		if nextOk {
			var defaultValue T
			return defaultValue, false
		}
	}

	// to reach here there were 0 or 1 items, ok is true
	return result, true
}

func LastItr[T any](src Iterator[T]) (T, bool) {
	itr := src.initItr()

	result, ok := itr.Next()

	if !ok {
		return result, ok
	}

	for {
		next, ok := itr.Next()
		if !ok {
			return result, true
		}
		result = next
	}
}

func LastOrDefaultItr[T any](src Iterator[T]) T {
	itr := src.initItr()

	result, ok := itr.Next()

	if !ok {
		return result
	}

	for {
		next, ok := itr.Next()
		if !ok {
			return result
		}
		result = next
	}
}
