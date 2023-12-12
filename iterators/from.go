package iterators

import cmn "github.com/alexmacinnes/golinq/common"

//Slice to Itr
type itrFromSlice[T any] struct {
	Input  *[]T
	index  int
	length int
}

func (x *itrFromSlice[T]) Next() (T, bool) {
	x.index += 1
	if x.index >= x.length {
		var none T
		return none, false
	}
	return (*x.Input)[x.index], true
}

type iteratorFromSlice[T any] struct {
	Input *[]T
}

func (x *iteratorFromSlice[T]) initItr() itr[T] {
	return &itrFromSlice[T]{
		Input:  x.Input,
		index:  -1,
		length: len(*x.Input),
	}
}

//Slice to PtrItr
type ptrItrFromSlice[T any] struct {
	Input  *[]T
	index  int
	length int
}

func (x *ptrItrFromSlice[T]) Next() (*T, bool) {
	x.index += 1
	if x.index >= x.length {
		return nil, false
	}
	return &((*x.Input)[x.index]), true
}

type ptrIteratorFromSlice[T any] struct {
	Input *[]T
}

func (x *ptrIteratorFromSlice[T]) initItr() itr[*T] {
	return &ptrItrFromSlice[T]{
		Input:  x.Input,
		index:  -1,
		length: len(*x.Input),
	}
}

type iteratorFromMap[T_Key comparable, T_Value any] struct {
	Input *map[T_Key]T_Value
}

func (x *iteratorFromMap[T_Key, T_Value]) initItr() itr[cmn.KeyValuePair[T_Key, T_Value]] {
	inputChannel := make(chan cmn.KeyValuePair[T_Key, T_Value])

	go func() {
		for k, v := range *(x.Input) {
			inputChannel <- cmn.KeyValuePair[T_Key, T_Value]{Key: k, Value: v}
		}
		close(inputChannel)
	}()

	return &itrFromChan[cmn.KeyValuePair[T_Key, T_Value]]{
		InputChannel: inputChannel,
	}
}

//Map to PtrItr
type ptrIteratorFromMap[T_Key comparable, T_Value any] struct {
	Input *map[T_Key]T_Value
}

func (x *ptrIteratorFromMap[T_Key, T_Value]) initItr() itr[cmn.KeyValuePair[T_Key, *T_Value]] {
	inputChannel := make(chan cmn.KeyValuePair[T_Key, *T_Value])

	go func() {
		for k, v := range *(x.Input) {
			inputChannel <- cmn.KeyValuePair[T_Key, *T_Value]{Key: k, Value: &v}
		}
		close(inputChannel)
	}()

	return &itrFromChan[cmn.KeyValuePair[T_Key, *T_Value]]{
		InputChannel: inputChannel,
	}
}

//Chan to Itr
type itrFromChan[T any] struct {
	InputChannel chan T
}

func (x *itrFromChan[T]) Next() (T, bool) {
	result, ok := <-x.InputChannel
	return result, ok
}

//Public
func FromSlice[T any](slice *[]T) Iterator[T] {
	return &iteratorFromSlice[T]{
		Input: slice,
	}
}

func PointersFromSlice[T any](slice *[]T) Iterator[*T] {
	return &ptrIteratorFromSlice[T]{
		Input: slice,
	}
}

func FromMap[T_Key comparable, T_Value any](input *map[T_Key]T_Value) Iterator[cmn.KeyValuePair[T_Key, T_Value]] {
	return &iteratorFromMap[T_Key, T_Value]{
		Input: input,
	}
}

func PointersFromMap[T_Key comparable, T_Value any](input *map[T_Key]T_Value) Iterator[cmn.KeyValuePair[T_Key, *T_Value]] {
	return &ptrIteratorFromMap[T_Key, T_Value]{
		Input: input,
	}
}
