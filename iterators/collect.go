package iterators

func Any[T any](src Iterator[T]) bool {
	itr := src.initItr()
	_, ok := itr.Next()
	return ok
}

func All[T any](src Iterator[T], predicate func(T) bool) bool {
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

func Contains[T comparable](src Iterator[T], item T) bool {
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

func ElementAt[T any](src Iterator[T], index int) (T, bool) {
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

func First[T any](src Iterator[T]) (T, bool) {
	itr := src.initItr()
	return itr.Next()
}

func FirstOrDefault[T any](src Iterator[T]) T {
	itr := src.initItr()
	result, _ := itr.Next()
	return result
}

func Single[T any](src Iterator[T]) (T, bool) {
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

func SingleOrDefault[T any](src Iterator[T]) (T, bool) {
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

func Last[T any](src Iterator[T]) (T, bool) {
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

func LastOrDefault[T any](src Iterator[T]) T {
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
