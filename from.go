package golinq

type KeyValuePair[T_Key comparable, T_Value any] struct {
	Key   T_Key
	Value T_Value
}

type enumerableFromSlice[T any] struct {
	Input *[]T
}

type ptrEnumerableFromSlice[T any] struct {
	Input *[]T
}

type enumerableFromMap[T_Key comparable, T_Value any] struct {
	Input *map[T_Key]T_Value
}

type ptrEnumerableFromMap[T_Key comparable, T_Value any] struct {
	Input *map[T_Key]T_Value
}

func (this *enumerableFromSlice[T_Out]) getAction() *actionDelegate[T_Out] {
	actionDelegate, ctx := newActionDelegate[T_Out]()

	action := func() {
		defer close(actionDelegate.ResultChannel)

		for _, x := range *this.Input {
			if actionIsCancelled(ctx) {
				break // abort the current operation
			}
			actionDelegate.ResultChannel <- x
		}
	}
	actionDelegate.Action = action

	return actionDelegate
}

func (this *ptrEnumerableFromSlice[T_Out]) getAction() *actionDelegate[*T_Out] {
	actionDelegate, ctx := newActionDelegate[*T_Out]()

	action := func() {
		defer close(actionDelegate.ResultChannel)

		for _, x := range *this.Input {
			if actionIsCancelled(ctx) {
				break // abort the current operation
			}
			actionDelegate.ResultChannel <- &x
		}
	}
	actionDelegate.Action = action

	return actionDelegate
}

func (this *enumerableFromMap[T_Key, T_Value]) getAction() *actionDelegate[KeyValuePair[T_Key, T_Value]] {
	actionDelegate, ctx := newActionDelegate[KeyValuePair[T_Key, T_Value]]()

	action := func() {
		defer close(actionDelegate.ResultChannel)

		for k, v := range *this.Input {
			if actionIsCancelled(ctx) {
				break // abort the current operation
			}
			kvp := KeyValuePair[T_Key, T_Value]{
				Key:   k,
				Value: v,
			}
			actionDelegate.ResultChannel <- kvp
		}
	}
	actionDelegate.Action = action

	return actionDelegate
}

func (this *ptrEnumerableFromMap[T_Key, T_Value]) getAction() *actionDelegate[KeyValuePair[T_Key, *T_Value]] {
	actionDelegate, ctx := newActionDelegate[KeyValuePair[T_Key, *T_Value]]()

	action := func() {
		defer close(actionDelegate.ResultChannel)

		for k, v := range *this.Input {
			if actionIsCancelled(ctx) {
				break // abort the current operation
			}
			kvp := KeyValuePair[T_Key, *T_Value]{
				Key:   k,
				Value: &v,
			}
			actionDelegate.ResultChannel <- kvp
		}
	}
	actionDelegate.Action = action

	return actionDelegate
}

func EnumerableFromSlice[T any](slice *[]T) Enumerable[T] {
	return &enumerableFromSlice[T]{Input: slice}
}

func PtrEnumerableFromSlice[T any](slice *[]T) Enumerable[*T] {
	return &ptrEnumerableFromSlice[T]{Input: slice}
}

func EnumerableFromMap[T_Key comparable, T_Value any](input *map[T_Key]T_Value) Enumerable[KeyValuePair[T_Key, T_Value]] {
	return &enumerableFromMap[T_Key, T_Value]{Input: input}
}

func PtrEnumerableFromMap[T_Key comparable, T_Value any](input *map[T_Key]T_Value) Enumerable[KeyValuePair[T_Key, *T_Value]] {
	//TODO - this does work consistently
	// it's probably wrong to address map vals
	// is this even useful?

	x := EnumerableFromMap(input)
	res := Select(x, func(kvp KeyValuePair[T_Key, T_Value]) KeyValuePair[T_Key, *T_Value] {
		return KeyValuePair[T_Key, *T_Value]{Key: kvp.Key, Value: &kvp.Value}
	})

	return res

	// intermittent failures
	//return &ptrEnumerableFromMap[T_Key, T_Value]{Input: input}
}
