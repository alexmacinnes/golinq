package enumerables

import "context"

type actionDelegate[T any] struct {
	Action        func()
	ResultChannel chan T
	CancelFunc    context.CancelFunc
}

type Enumerable[T any] interface {
	getAction() *actionDelegate[T]
}

func newActionDelegate[T any]() (*actionDelegate[T], *context.Context) {
	ctx, cancelFunc := context.WithCancel(context.Background())

	actionDelegate := actionDelegate[T]{
		Action:        nil,
		ResultChannel: make(chan T),
		CancelFunc:    cancelFunc,
	}

	return &actionDelegate, &ctx
}

func actionIsCancelled(ctx *context.Context) bool {
	select {
	case <-(*ctx).Done():
		return true
	default:
		return false
	}
}

func runAction[T any](src Enumerable[T]) (chan T, func()) {
	actionDelegate := src.getAction()

	go actionDelegate.Action()
	return actionDelegate.ResultChannel, actionDelegate.CancelFunc
}
