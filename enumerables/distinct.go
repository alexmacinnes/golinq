package enumerables

type enumerableDistinct[T comparable] struct {
	Prior Enumerable[T]
}

func (this *enumerableDistinct[T]) getAction() *actionDelegate[T] {
	actionDelegate, ctx := newActionDelegate[T]()

	action := func() {
		defer close(actionDelegate.ResultChannel)

		previousItems := make(map[T]bool)

		priorAction := this.Prior.getAction()

		go priorAction.Action()
		chanIn := priorAction.ResultChannel

		for x := range chanIn {
			if actionIsCancelled(ctx) {
				priorAction.CancelFunc() // cancel the prior operation
				break
			}
			if !previousItems[x] {
				previousItems[x] = true
				actionDelegate.ResultChannel <- x
			}
		}
	}
	actionDelegate.Action = action

	return actionDelegate
}

func Distinct[T comparable](prior Enumerable[T]) Enumerable[T] {
	return &enumerableDistinct[T]{
		Prior: prior,
	}
}
