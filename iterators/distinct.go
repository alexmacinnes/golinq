package iterators

type itrDistinct[T comparable] struct {
	Inner         itr[T]
	previousItems map[T]interface{}
}

func (x *itrDistinct[T]) Next() (T, bool) {
	for {
		next, ok := x.Inner.Next()
		if !ok {
			return next, false
		}

		_, exists := x.previousItems[next]
		if !exists {
			x.previousItems[next] = true
			return next, true
		}
	}
}

type iteratorDistinct[T comparable] struct {
	Inner Iterator[T]
}

func (x *iteratorDistinct[T]) initItr() itr[T] {
	return &itrDistinct[T]{
		Inner:         x.Inner.initItr(),
		previousItems: make(map[T]interface{}),
	}
}

func DistinctItr[T comparable](inner Iterator[T]) Iterator[T] {
	return &iteratorDistinct[T]{
		Inner: inner,
	}
}
