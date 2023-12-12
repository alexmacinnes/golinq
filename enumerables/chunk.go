package enumerables

type enumerableChunk[T any] struct {
	Prior     Enumerable[T]
	ChunkSize uint32
}

func (this *enumerableChunk[T]) getAction() *actionDelegate[[]T] {
	actionDelegate, ctx := newActionDelegate[[]T]()

	action := func() {
		defer close(actionDelegate.ResultChannel)

		priorAction := this.Prior.getAction()

		go priorAction.Action()
		chanIn := priorAction.ResultChannel

		currentChunk := []T{}
		currentCount := 0

		for x := range chanIn {
			if actionIsCancelled(ctx) {
				priorAction.CancelFunc() // cancel the prior operation
				break
			}

			currentChunk = append(currentChunk, x)
			currentCount++

			if currentCount == int(this.ChunkSize) {
				actionDelegate.ResultChannel <- currentChunk
				currentChunk = []T{}
				currentCount = 0
			}
		}

		if currentCount > 0 {
			actionDelegate.ResultChannel <- currentChunk
		}

	}
	actionDelegate.Action = action

	return actionDelegate
}

func Chunk[T any](prior Enumerable[T], chunkSize uint32) Enumerable[[]T] {
	return &enumerableChunk[T]{
		Prior:     prior,
		ChunkSize: chunkSize,
	}
}
