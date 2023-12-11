package golinq

type itrChunk[T any] struct {
	Inner             itr[T]
	ChunkSize         uint32
	CurrentChunkCount uint32
	CurrentChunk      []T
}

func (x *itrChunk[T]) Next() ([]T, bool) {
	for {
		next, ok := x.Inner.Next()
		if !ok {
			if x.CurrentChunkCount > 0 {
				x.CurrentChunkCount = 0
				return x.CurrentChunk, true
			} else {
				var none []T
				return none, false
			}
		}

		x.CurrentChunk = append(x.CurrentChunk, next)
		x.CurrentChunkCount++

		if x.CurrentChunkCount == x.ChunkSize {
			res := x.CurrentChunk
			x.CurrentChunkCount = 0
			x.CurrentChunk = []T{}
			return res, true
		}
	}
}

type iteratorChunk[T any] struct {
	Inner     Iterator[T]
	ChunkSize uint32
}

func (x *iteratorChunk[T]) initItr() itr[[]T] {
	return &itrChunk[T]{
		Inner:             x.Inner.initItr(),
		ChunkSize:         x.ChunkSize,
		CurrentChunkCount: 0,
		CurrentChunk:      []T{},
	}
}

func ChunkItr[T any](inner Iterator[T], chunkSize uint32) Iterator[[]T] {
	return &iteratorChunk[T]{
		Inner:     inner,
		ChunkSize: chunkSize,
	}
}
