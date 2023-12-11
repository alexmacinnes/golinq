package golinq

type itr[T any] interface {
	Next() (T, bool)
}

type Iterator[T any] interface {
	initItr() itr[T]
}
