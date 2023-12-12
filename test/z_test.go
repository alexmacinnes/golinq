package test

import (
	cmn "github.com/alexmacinnes/golinq/common"
	enm "github.com/alexmacinnes/golinq/enumerables"
	itr "github.com/alexmacinnes/golinq/iterators"

	"fmt"
	"reflect"
	"strings"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

func personSlice0() []Person {
	return []Person{}
}

func personSlice1() []Person {
	return []Person{
		{"James", 23},
	}
}

func personSlice5() []Person {
	return []Person{
		{"James", 23},
		{"Lucy", 33},
		{"Zack", 41},
		{"Abi", 19},
		{"Rach", 33},
	}
}

func personMap0() map[string]Person {
	return map[string]Person{}
}

func personMap5() map[string]Person {
	result := map[string]Person{}
	for _, x := range personSlice5() {
		result[x.Name] = x
	}
	return result
}

func assertResult[T any](t *testing.T, expected T, actual T) {
	if !reflect.DeepEqual(expected, actual) {
		msg := fmt.Sprintf("Expected %v, Actual %v", expected, actual)
		(*t).Fatal(msg)
	}
}

func personAge(p Person) int { return p.Age }

func personName(p Person) string { return p.Name }

func intRange(min int, max int) []int {
	nums := []int{}
	for i := min; i <= max; i++ {
		nums = append(nums, i)
	}
	return nums
}

func TestSelect(t *testing.T) {
	p := personSlice1()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, func(p Person) int { return p.Age })
	x3 := enm.Select(x2, func(i int) int { return 0 - i })

	actual1 := enm.EnumerableToSlice(x3)
	expected1 := []int{-23}
	assertResult(t, expected1, actual1)

	p = append(p, Person{"Zane", 20})

	actual2 := enm.EnumerableToSlice(x3)
	expected2 := []int{-23, -20}
	assertResult(t, expected2, actual2)
}

func TestSelect_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, func(p Person) int { return p.Age })
	x3 := itr.SelectItr(x2, func(i int) int { return 0 - i })

	actual1 := itr.IteratorToSlice(x3)
	expected1 := []int{-23}
	assertResult(t, expected1, actual1)

	p = append(p, Person{"Zane", 20})

	actual2 := itr.IteratorToSlice(x3)
	expected2 := []int{-23, -20}
	assertResult(t, expected2, actual2)
}

func TestWhereEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Where(x1, func(p Person) bool { return p.Age > 40 })
	x3 := enm.Select(x2, personName)
	where := enm.EnumerableToSlice(x3)

	assertResult(t, []string{}, where)
}

func TestWhereEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.WhereItr(x1, func(p Person) bool { return p.Age > 40 })
	x3 := itr.SelectItr(x2, personName)
	where := itr.IteratorToSlice(x3)

	assertResult(t, []string{}, where)
}

func TestWhere(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Where(x1, func(p Person) bool { return p.Age > 30 })
	x3 := enm.Select(x2, personName)
	where := enm.EnumerableToSlice(x3)

	assertResult(t, []string{"Lucy", "Zack", "Rach"}, where)
}

func TestWhere_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.WhereItr(x1, func(p Person) bool { return p.Age > 30 })
	x3 := itr.SelectItr(x2, personName)
	where := itr.IteratorToSlice(x3)

	assertResult(t, []string{"Lucy", "Zack", "Rach"}, where)
}

func TestAnyFalse(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	actual := enm.Any(x1)

	assertResult(t, false, actual)
}

func TestAnyFalse_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	actual := itr.AnyItr(x1)

	assertResult(t, false, actual)
}

func TestAnyTrue(t *testing.T) {
	nums := intRange(1, 10000)

	numsProcessed := 0
	x1 := enm.EnumerableFromSlice(&nums)
	x2 := enm.Select(x1, func(x int) int {
		numsProcessed++
		return x
	})
	any := enm.Any(x2)

	assertResult(t, true, any)

	// should stop processing as soon as it finds first item
	if numsProcessed < 1 || numsProcessed > 2 {
		panic("Expected count 1 or 2, got " + string(rune(numsProcessed)))
	}
}

func TestAnyTrue_Itr(t *testing.T) {
	nums := intRange(1, 10000)

	numsProcessed := 0
	x1 := itr.IteratorFromSlice(&nums)
	x2 := itr.SelectItr(x1, func(x int) int {
		numsProcessed++
		return x
	})
	any := itr.AnyItr(x2)

	assertResult(t, true, any)

	// should stop processing as soon as it finds first item
	if numsProcessed != 1 {
		panic("Expected count 1, got " + string(rune(numsProcessed)))
	}
}

func TestAllEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	actual := enm.All(x1, func(p Person) bool { return p.Age > 1000 })

	assertResult(t, true, actual) // always true on empty input
}

func TestAllEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	actual := itr.AllItr(x1, func(p Person) bool { return p.Age > 1000 })

	assertResult(t, true, actual) // always true on empty input
}

func TestAllTrue(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	actual := enm.All(x1, func(p Person) bool { return p.Age > 0 })

	assertResult(t, true, actual)
}

func TestAllTrue_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	actual := itr.AllItr(x1, func(p Person) bool { return p.Age > 0 })

	assertResult(t, true, actual)
}

func TestAllFalse(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	actual := enm.All(x1, func(p Person) bool { return p.Age > 21 })

	assertResult(t, false, actual)
}

func TestAllFalse_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	actual := itr.AllItr(x1, func(p Person) bool { return p.Age > 21 })

	assertResult(t, false, actual)
}

func TestMaxEmpty(t *testing.T) {
	nums := personSlice0()

	x1 := enm.EnumerableFromSlice(&nums)
	x2 := enm.Select(x1, personAge)
	_, ok := enm.Max(x2)

	assertResult(t, false, ok)
}

func TestMaxEmpty_Itr(t *testing.T) {
	nums := personSlice0()

	x1 := itr.IteratorFromSlice(&nums)
	x2 := itr.SelectItr(x1, personAge)
	_, ok := itr.MaxItr(x2)

	assertResult(t, false, ok)
}

func TestMaxOneItem(t *testing.T) {
	nums := personSlice1()

	x1 := enm.EnumerableFromSlice(&nums)
	x2 := enm.Select(x1, personAge)
	max, ok := enm.Max(x2)

	assertResult(t, 23, max)
	assertResult(t, true, ok)
}

func TestMaxOneItem_Itr(t *testing.T) {
	nums := personSlice1()

	x1 := itr.IteratorFromSlice(&nums)
	x2 := itr.SelectItr(x1, personAge)
	max, ok := itr.MaxItr(x2)

	assertResult(t, 23, max)
	assertResult(t, true, ok)
}

func TestMaxString(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	max, ok := enm.Max(x2)

	assertResult(t, "Zack", max)
	assertResult(t, true, ok)
}

func TestMaxString_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	max, ok := itr.MaxItr(x2)

	assertResult(t, "Zack", max)
	assertResult(t, true, ok)
}

func TestMinEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	_, ok := enm.Min(x2)

	assertResult(t, false, ok)
}

func TestMinEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	_, ok := itr.MinItr(x2)

	assertResult(t, false, ok)
}

func TestMinOneItem(t *testing.T) {
	p := personSlice1()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	min, ok := enm.Min(x2)

	assertResult(t, 23, min)
	assertResult(t, true, ok)
}

func TestMinOneItem_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	min, ok := itr.MinItr(x2)

	assertResult(t, 23, min)
	assertResult(t, true, ok)
}

func TestMinString(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	min, ok := enm.Min(x2)

	assertResult(t, "Abi", min)
	assertResult(t, true, ok)
}

func TestMinString_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	min, ok := itr.MinItr(x2)

	assertResult(t, "Abi", min)
	assertResult(t, true, ok)
}

func TestAvgEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	_, ok := enm.Avg(x2)

	assertResult(t, false, ok)
}

func TestAvgEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	_, ok := itr.AvgItr(x2)

	assertResult(t, false, ok)
}

func TestAvgOneItem(t *testing.T) {
	p := personSlice1()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	avg, ok := enm.Avg(x2)

	assertResult(t, 23.0, avg)
	assertResult(t, true, ok)
}

func TestAvgOneItem_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	avg, ok := itr.AvgItr(x2)

	assertResult(t, 23.0, avg)
	assertResult(t, true, ok)
}

func TestAvgManyItems(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	avg, ok := enm.Avg(x2)

	assertResult(t, 29.8, avg)
	assertResult(t, true, ok)
}

func TestAvgManyItems_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	avg, ok := itr.AvgItr(x2)

	assertResult(t, 29.8, avg)
	assertResult(t, true, ok)
}

func TestCountEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	count := enm.Count(x2)

	assertResult(t, 0, count)
}

func TestCountEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	count := itr.CountItr(x2)

	assertResult(t, 0, count)
}

func TestCountOneItem(t *testing.T) {
	p := personSlice1()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	count := enm.Count(x2)

	assertResult(t, 1, count)
}

func TestCountOneItem_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	count := itr.CountItr(x2)

	assertResult(t, 1, count)
}

func TestCount5(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	count := enm.Count(x2)

	assertResult(t, 5, count)
}

func TestCount5_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	count := itr.CountItr(x2)

	assertResult(t, 5, count)
}

func TestSumEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	count := enm.Count(x2)

	assertResult(t, 0, count)
}

func TestSumEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	count := itr.CountItr(x2)

	assertResult(t, 0, count)
}

func TestSumOneItem(t *testing.T) {
	p := personSlice1()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	sum := enm.Sum(x2)

	assertResult(t, 23, sum)
}

func TestSumOneItem_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	sum := itr.SumItr(x2)

	assertResult(t, 23, sum)
}

func TestSum5(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	sum := enm.Sum(x2)

	assertResult(t, 149, sum)
}

func TestSum5_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	sum := itr.SumItr(x2)

	assertResult(t, 149, sum)
}

func TestSumFiltered(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Where(x1, func(p Person) bool { return p.Age > 30 })
	x3 := enm.Select(x2, personAge)
	sum := enm.Sum(x3)

	assertResult(t, 107, sum)
}

func TestSumFiltered_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.WhereItr(x1, func(p Person) bool { return p.Age > 30 })
	x3 := itr.SelectItr(x2, personAge)
	sum := itr.SumItr(x3)

	assertResult(t, 107, sum)
}

func TestAccumulateEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	accum := enm.Accumulate(x1, 99, func(prior int, p Person) int { return prior + p.Age })

	assertResult(t, 99, accum)
}

func TestAccumulateEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	accum := itr.AccumulateItr(x1, 99, func(prior int, p Person) int { return prior + p.Age })

	assertResult(t, 99, accum)
}

func TestAccumulateManyInts(t *testing.T) {
	nums := intRange(2, 4)

	x1 := enm.EnumerableFromSlice(&nums)
	accum := enm.Accumulate(x1, -100, func(prior int, item int) int { return prior * item })

	assertResult(t, -2400, accum) // -100 * 2 * 3 * 4
}

func TestAccumulateManyInts_Itr(t *testing.T) {
	nums := intRange(2, 4)

	x1 := itr.IteratorFromSlice(&nums)
	accum := itr.AccumulateItr(x1, -100, func(prior int, item int) int { return prior * item })

	assertResult(t, -2400, accum) // -100 * 2 * 3 * 4
}

func TestAccumulateManyStrings(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Where(x1, func(p Person) bool { return p.Age < 30 })
	accum := enm.Accumulate(x2, "PeopleUnder30:", func(prior string, p Person) string { return prior + p.Name })

	assertResult(t, "PeopleUnder30:JamesAbi", accum) // -100 * 2 * 3 * 4
}

func TestAccumulateManyStrings_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.WhereItr(x1, func(p Person) bool { return p.Age < 30 })
	accum := itr.AccumulateItr(x2, "PeopleUnder30:", func(prior string, p Person) string { return prior + p.Name })

	assertResult(t, "PeopleUnder30:JamesAbi", accum) // -100 * 2 * 3 * 4
}

func TestContainsEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	contains := enm.Contains(x2, 99)

	assertResult(t, false, contains)
}

func TestContainsEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	contains := itr.ContainsItr(x2, 99)

	assertResult(t, false, contains)
}

func TestContainsMany(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	contains999 := enm.Contains(x2, 999)
	contains23 := enm.Contains(x2, 23)

	assertResult(t, false, contains999)
	assertResult(t, true, contains23)
}

func TestContainsMany_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	contains999 := itr.ContainsItr(x2, 999)
	contains23 := itr.ContainsItr(x2, 23)

	assertResult(t, false, contains999)
	assertResult(t, true, contains23)
}

func TestElementAt(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	element0, ok0 := enm.ElementAt(x2, 0)
	element2, ok2 := enm.ElementAt(x2, 2)
	element4, ok4 := enm.ElementAt(x2, 4)
	element5, ok5 := enm.ElementAt(x2, 5)

	assertResult(t, true, ok0)
	assertResult(t, true, ok2)
	assertResult(t, true, ok4)
	assertResult(t, false, ok5)

	assertResult(t, "James", element0)
	assertResult(t, "Zack", element2)
	assertResult(t, "Rach", element4)
	assertResult(t, "", element5)
}

func TestElementAt_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	element0, ok0 := itr.ElementAtItr(x2, 0)
	element2, ok2 := itr.ElementAtItr(x2, 2)
	element4, ok4 := itr.ElementAtItr(x2, 4)
	element5, ok5 := itr.ElementAtItr(x2, 5)

	assertResult(t, true, ok0)
	assertResult(t, true, ok2)
	assertResult(t, true, ok4)
	assertResult(t, false, ok5)

	assertResult(t, "James", element0)
	assertResult(t, "Zack", element2)
	assertResult(t, "Rach", element4)
	assertResult(t, "", element5)
}

func TestFirstEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.First(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestFirstEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first, ok := itr.FirstItr(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestFirst(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.First(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestFirst_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first, ok := itr.FirstItr(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestFirstOrDefaultEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first := enm.FirstOrDefault(x2)

	assertResult(t, "", first)
}

func TestFirstOrDefaultEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first := itr.FirstOrDefaultItr(x2)

	assertResult(t, "", first)
}

func TestFirstOrDefault(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first := enm.FirstOrDefault(x2)

	assertResult(t, "James", first)
}

func TestFirstOrDefault_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first := itr.FirstOrDefaultItr(x2)

	assertResult(t, "James", first)
}

func TestSingleEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingleEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first, ok := itr.SingleItr(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingle(t *testing.T) {
	p := personSlice1()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Single(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestSingle_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first, ok := itr.SingleItr(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestSingleMany(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingleMany_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first, ok := itr.SingleItr(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingleOrDefaultEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.SingleOrDefault(x2)

	assertResult(t, "", first)
	assertResult(t, true, ok)
}

func TestSingleOrDefaultEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first, ok := itr.SingleOrDefaultItr(x2)

	assertResult(t, "", first)
	assertResult(t, true, ok)
}

func TestSingleOrDefault(t *testing.T) {
	p := personSlice1()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.SingleOrDefault(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestSingleOrDefault_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first, ok := itr.SingleOrDefaultItr(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestSingleOrDefaultMany(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingleOrDefaultMany_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	first, ok := itr.SingleItr(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestLastEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	last, ok := enm.Last(x2)

	assertResult(t, "", last)
	assertResult(t, false, ok)
}

func TestLastEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	last, ok := itr.LastItr(x2)

	assertResult(t, "", last)
	assertResult(t, false, ok)
}

func TestLast(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Last(x2)

	assertResult(t, "Rach", first)
	assertResult(t, true, ok)
}

func TestLast_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	last, ok := itr.LastItr(x2)

	assertResult(t, "Rach", last)
	assertResult(t, true, ok)
}

func TestLastOrDefaultEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	last := enm.LastOrDefault(x2)

	assertResult(t, "", last)
}

func TestLastOrDefaultEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	last := itr.LastOrDefaultItr(x2)

	assertResult(t, "", last)
}

func TestLastOrDefault(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	last := enm.LastOrDefault(x2)

	assertResult(t, "Rach", last)
}

func TestLastOrDefault_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	last := itr.LastOrDefaultItr(x2)

	assertResult(t, "Rach", last)
}

func TestChunkEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	x3 := enm.Chunk(x2, 3)
	chunks := enm.EnumerableToSlice(x3)

	assertResult(t, [][]string{}, chunks)
}

func TestChunkEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	x3 := itr.ChunkItr(x2, 3)
	chunks := itr.IteratorToSlice(x3)

	assertResult(t, [][]string{}, chunks)
}

func TestChunk(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personName)
	x3 := enm.Chunk(x2, 3)
	chunks := enm.EnumerableToSlice(x3)

	expected := [][]string{
		{"James", "Lucy", "Zack"},
		{"Abi", "Rach"},
	}
	assertResult(t, expected, chunks)
}

func TestChunk_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personName)
	x3 := itr.ChunkItr(x2, 3)
	chunks := itr.IteratorToSlice(x3)

	expected := [][]string{
		{"James", "Lucy", "Zack"},
		{"Abi", "Rach"},
	}
	assertResult(t, expected, chunks)
}

func TestDistinctEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	x3 := enm.Distinct(x2)
	distinct := enm.EnumerableToSlice(x3)

	assertResult(t, []int{}, distinct)
}

func TestDistinctEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	x3 := itr.DistinctItr(x2)
	distinct := itr.IteratorToSlice(x3)

	assertResult(t, []int{}, distinct)
}

func TestDistinct(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	x2 := enm.Select(x1, personAge)
	x3 := enm.Distinct(x2)
	distinct := enm.EnumerableToSlice(x3)

	assertResult(t, []int{23, 33, 41, 19}, distinct)
}

func TestDistinct_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, personAge)
	x3 := itr.DistinctItr(x2)
	distinct := itr.IteratorToSlice(x3)

	assertResult(t, []int{23, 33, 41, 19}, distinct)
}

func TestEnumerableToMap(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	mapped, ok := enm.EnumerableToMap(x1,
		func(p Person) string { return strings.ToUpper(p.Name) },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, true)
	assertResult(t, mapped, map[string]int{"JAMES": 230, "LUCY": 330, "ZACK": 410, "ABI": 190, "RACH": 330})
}

func TestIteratorToMap(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	mapped, ok := itr.IteratorToMap(x1,
		func(p Person) string { return strings.ToUpper(p.Name) },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, true)
	assertResult(t, mapped, map[string]int{"JAMES": 230, "LUCY": 330, "ZACK": 410, "ABI": 190, "RACH": 330})
}

func TestEnumerableToDictionaryEmpty(t *testing.T) {
	p := personSlice0()

	x1 := enm.EnumerableFromSlice(&p)
	mapped, ok := enm.EnumerableToMap(x1,
		func(p Person) string { return strings.ToUpper(p.Name) },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, true)
	assertResult(t, len(mapped), 0)
}

func TestIteratorToDictionaryEmpty(t *testing.T) {
	p := personSlice0()

	x1 := itr.IteratorFromSlice(&p)
	mapped, ok := itr.IteratorToMap(x1,
		func(p Person) string { return strings.ToUpper(p.Name) },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, true)
	assertResult(t, len(mapped), 0)
}

func TestEnumerableToDictionaryDuplicateKey(t *testing.T) {
	p := personSlice5()

	x1 := enm.EnumerableFromSlice(&p)
	mapped, ok := enm.EnumerableToMap(x1,
		func(p Person) string { return "DUPLICATE" },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, false)
	assertResult(t, mapped, nil)
}

func TestIteratorToDictionaryDuplicateKey(t *testing.T) {
	p := personSlice5()

	x1 := itr.IteratorFromSlice(&p)
	mapped, ok := itr.IteratorToMap(x1,
		func(p Person) string { return "DUPLICATE" },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, false)
	assertResult(t, mapped, nil)
}

func TestPtrEnumerableFromSlice(t *testing.T) {
	p := personSlice5()

	x1 := enm.PtrEnumerableFromSlice(&p)
	x2 := enm.Select(x1, func(ptr *Person) int { return (*ptr).Age })
	maxAge, _ := enm.Max(x2)

	assertResult(t, maxAge, 41)
}

func TestPtrIteratorFromSlice(t *testing.T) {
	p := personSlice5()

	x1 := itr.PtrIteratorFromSlice(&p)
	x2 := itr.SelectItr(x1, func(ptr *Person) int { return (*ptr).Age })
	maxAge, _ := itr.MaxItr(x2)

	assertResult(t, maxAge, 41)
}

func TestEnumerableFromMap(t *testing.T) {
	p := personMap5()

	x1 := enm.EnumerableFromMap(&p)
	x2 := enm.Where(x1, func(kvp cmn.KeyValuePair[string, Person]) bool { return kvp.Key == "Abi" })
	x3 := enm.Select(x2, func(kvp cmn.KeyValuePair[string, Person]) int { return kvp.Value.Age })
	abiAge, _ := enm.Single(x3)

	assertResult(t, abiAge, 19)
}

func TestIteratorFromMap(t *testing.T) {
	p := personMap5()

	x1 := itr.IteratorFromMap(&p)
	x2 := itr.WhereItr(x1, func(kvp cmn.KeyValuePair[string, Person]) bool { return kvp.Key == "Abi" })
	x3 := itr.SelectItr(x2, func(kvp cmn.KeyValuePair[string, Person]) int { return kvp.Value.Age })
	abiAge, _ := itr.SingleItr(x3)

	assertResult(t, abiAge, 19)
}

func TestEnumerableFromMapEmpty(t *testing.T) {
	p := personMap0()

	x1 := enm.EnumerableFromMap(&p)
	slice := enm.EnumerableToSlice(x1)

	assertResult(t, len(slice), 0)
}

func TestIteratorFromMapEmpty(t *testing.T) {
	p := personMap0()

	x1 := itr.IteratorFromMap(&p)
	slice := itr.IteratorToSlice(x1)

	assertResult(t, len(slice), 0)
}

func TestPtrEnumerableFromMap(t *testing.T) {
	//TODO - intermittent failures
	p := personMap5()

	x1 := enm.PtrEnumerableFromMap(&p)
	x2 := enm.Where(x1, func(kvp cmn.KeyValuePair[string, *Person]) bool { return kvp.Key == "Abi" })
	x3 := enm.Select(x2, func(kvp cmn.KeyValuePair[string, *Person]) int { return (*kvp.Value).Age })
	abiAge, _ := enm.Single(x3)

	assertResult(t, abiAge, 19)
}

func TestPtrIteratorFromMap(t *testing.T) {
	//TODO - intermittent failures
	p := personMap5()

	x1 := itr.PtrIteratorFromMap(&p)
	x2 := itr.WhereItr(x1, func(kvp cmn.KeyValuePair[string, *Person]) bool { return kvp.Key == "Abi" })
	x3 := itr.SelectItr(x2, func(kvp cmn.KeyValuePair[string, *Person]) int { return (*kvp.Value).Age })
	abiAge, _ := itr.SingleItr(x3)

	assertResult(t, abiAge, 19)
}

func TestPtrEnumerableFromMapEmpty(t *testing.T) {
	p := personMap0()

	x1 := enm.PtrEnumerableFromMap(&p)
	slice := enm.EnumerableToSlice(x1)

	assertResult(t, len(slice), 0)
}

func TestPtrIteratorFromMapEmpty(t *testing.T) {
	p := personMap0()

	x1 := itr.PtrIteratorFromMap(&p)
	slice := itr.IteratorToSlice(x1)

	assertResult(t, len(slice), 0)
}
