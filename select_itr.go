package golinq

type itrSelect[T_In any, T_Out any] struct {
	Inner    itr[T_In]
	Selector func(T_In) T_Out
}

func (x *itrSelect[T_In, T_Out]) Next() (T_Out, bool) {
	next, ok := x.Inner.Next()
	if !ok {
		var none T_Out
		return none, false
	}

	return x.Selector(next), true
}

type iteratorSelect[T_In any, T_Out any] struct {
	Inner    Iterator[T_In]
	Selector func(T_In) T_Out
}

func (x *iteratorSelect[T_In, T_Out]) initItr() itr[T_Out] {
	return &itrSelect[T_In, T_Out]{
		Inner:    x.Inner.initItr(),
		Selector: x.Selector,
	}
}

func SelectItr[T_In any, T_Out any](inner Iterator[T_In], selector func(T_In) T_Out) Iterator[T_Out] {
	return &iteratorSelect[T_In, T_Out]{
		Inner:    inner,
		Selector: selector,
	}
}
