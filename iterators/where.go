package iterators

type itrWhere[T any] struct {
	Inner     itr[T]
	Predicate func(T) bool
}

func (x *itrWhere[T]) Next() (T, bool) {
	for {
		next, ok := x.Inner.Next()
		if !ok {
			return next, false
		}

		match := x.Predicate(next)
		if match {
			return next, true
		}
	}
}

type iteratorWhere[T any] struct {
	Inner     Iterator[T]
	Predicate func(T) bool
}

func (x *iteratorWhere[T]) initItr() itr[T] {
	return &itrWhere[T]{
		Inner:     x.Inner.initItr(),
		Predicate: x.Predicate,
	}
}

func Where[T any](inner Iterator[T], predicate func(T) bool) Iterator[T] {
	return &iteratorWhere[T]{
		Inner:     inner,
		Predicate: predicate,
	}
}
