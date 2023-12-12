package enumerables

type enumerableSelect[T_In any, T_Out any] struct {
	Prior    Enumerable[T_In]
	Selector func(T_In) T_Out
}

func (this *enumerableSelect[T_In, T_Out]) getAction() *actionDelegate[T_Out] {
	actionDelegate, ctx := newActionDelegate[T_Out]()

	action := func() {
		defer close(actionDelegate.ResultChannel)

		priorAction := this.Prior.getAction()

		go priorAction.Action()
		chanIn := priorAction.ResultChannel

		for x := range chanIn {
			if actionIsCancelled(ctx) {
				priorAction.CancelFunc() // cancel the prior operation
				break
			}
			converted := this.Selector(x)
			actionDelegate.ResultChannel <- converted
		}
	}
	actionDelegate.Action = action

	return actionDelegate
}

func Select[T_In any, T_Out any](prior Enumerable[T_In], selector func(T_In) T_Out) Enumerable[T_Out] {
	return &enumerableSelect[T_In, T_Out]{
		Prior:    prior,
		Selector: selector,
	}
}
