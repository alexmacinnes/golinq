package iterators

import cmn "github.com/alexmacinnes/golinq/common"

func MaxItr[T cmn.Ordered](src Iterator[T]) (T, bool) {
	itr := src.initItr()

	max, ok := itr.Next()
	if !ok {
		return max, false
	}

	for {
		next, ok := itr.Next()
		if !ok {
			return max, true
		}
		if max < next {
			max = next
		}
	}
}

func MinItr[T cmn.Ordered](src Iterator[T]) (T, bool) {
	itr := src.initItr()

	min, ok := itr.Next()
	if !ok {
		return min, false
	}

	for {
		next, ok := itr.Next()
		if !ok {
			return min, true
		}
		if min > next {
			min = next
		}
	}
}

func AvgItr[T cmn.Numeric](src Iterator[T]) (float64, bool) {
	itr := src.initItr()

	var total float64 = 0
	count := 0

	for {
		next, ok := itr.Next()
		if !ok {
			break
		}
		total += float64(next)
		count++
	}

	if count == 0 {
		return 0, false
	}
	return total / float64(count), true
}

func SumItr[T cmn.Numeric](src Iterator[T]) T {
	itr := src.initItr()

	var total T = 0

	for {
		next, ok := itr.Next()
		if !ok {
			return total
		}
		total += next
	}
}

func CountItr[T any](src Iterator[T]) uint32 {
	itr := src.initItr()

	count := uint32(0)

	for {
		_, ok := itr.Next()
		if !ok {
			return count
		}
		count++
	}
}

func AccumulateItr[TAccumulate any, TItem any](src Iterator[TItem], seed TAccumulate, accumulator func(TAccumulate, TItem) TAccumulate) TAccumulate {
	itr := src.initItr()

	result := seed

	for {
		next, ok := itr.Next()
		if !ok {
			return result
		}
		result = accumulator(result, next)
	}
}
