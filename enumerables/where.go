package enumerables

type enumerableWhere[T any] struct {
	Prior     Enumerable[T]
	Predicate func(T) bool
}

func (this *enumerableWhere[T]) getAction() *actionDelegate[T] {
	actionDelegate, ctx := newActionDelegate[T]()

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
			if this.Predicate(x) {
				actionDelegate.ResultChannel <- x
			}
		}
	}
	actionDelegate.Action = action

	return actionDelegate
}

func Where[T any](prior Enumerable[T], predicate func(T) bool) Enumerable[T] {
	return &enumerableWhere[T]{
		Prior:     prior,
		Predicate: predicate,
	}
}
