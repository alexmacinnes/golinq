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

func TestSelect_Enm(t *testing.T) {
	p := personSlice1()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, func(p Person) int { return p.Age })
	x3 := enm.Select(x2, func(i int) int { return 0 - i })

	actual1 := enm.ToSlice(x3)
	expected1 := []int{-23}
	assertResult(t, expected1, actual1)

	p = append(p, Person{"Zane", 20})

	actual2 := enm.ToSlice(x3)
	expected2 := []int{-23, -20}
	assertResult(t, expected2, actual2)
}

func TestSelect_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, func(p Person) int { return p.Age })
	x3 := itr.Select(x2, func(i int) int { return 0 - i })

	actual1 := itr.ToSlice(x3)
	expected1 := []int{-23}
	assertResult(t, expected1, actual1)

	p = append(p, Person{"Zane", 20})

	actual2 := itr.ToSlice(x3)
	expected2 := []int{-23, -20}
	assertResult(t, expected2, actual2)
}

func TestWhereEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Where(x1, func(p Person) bool { return p.Age > 40 })
	x3 := enm.Select(x2, personName)
	where := enm.ToSlice(x3)

	assertResult(t, []string{}, where)
}

func TestWhereEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Where(x1, func(p Person) bool { return p.Age > 40 })
	x3 := itr.Select(x2, personName)
	where := itr.ToSlice(x3)

	assertResult(t, []string{}, where)
}

func TestWhere_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Where(x1, func(p Person) bool { return p.Age > 30 })
	x3 := enm.Select(x2, personName)
	where := enm.ToSlice(x3)

	assertResult(t, []string{"Lucy", "Zack", "Rach"}, where)
}

func TestWhere_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Where(x1, func(p Person) bool { return p.Age > 30 })
	x3 := itr.Select(x2, personName)
	where := itr.ToSlice(x3)

	assertResult(t, []string{"Lucy", "Zack", "Rach"}, where)
}

func TestAnyFalse_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	actual := enm.Any(x1)

	assertResult(t, false, actual)
}

func TestAnyFalse_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	actual := itr.Any(x1)

	assertResult(t, false, actual)
}

func TestAnyTrue_Enm(t *testing.T) {
	nums := intRange(1, 10000)

	numsProcessed := 0
	x1 := enm.FromSlice(&nums)
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
	x1 := itr.FromSlice(&nums)
	x2 := itr.Select(x1, func(x int) int {
		numsProcessed++
		return x
	})
	any := itr.Any(x2)

	assertResult(t, true, any)

	// should stop processing as soon as it finds first item
	if numsProcessed != 1 {
		panic("Expected count 1, got " + string(rune(numsProcessed)))
	}
}

func TestAllEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	actual := enm.All(x1, func(p Person) bool { return p.Age > 1000 })

	assertResult(t, true, actual) // always true on empty input
}

func TestAllEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	actual := itr.All(x1, func(p Person) bool { return p.Age > 1000 })

	assertResult(t, true, actual) // always true on empty input
}

func TestAllTrue_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	actual := enm.All(x1, func(p Person) bool { return p.Age > 0 })

	assertResult(t, true, actual)
}

func TestAllTrue_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	actual := itr.All(x1, func(p Person) bool { return p.Age > 0 })

	assertResult(t, true, actual)
}

func TestAllFalse_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	actual := enm.All(x1, func(p Person) bool { return p.Age > 21 })

	assertResult(t, false, actual)
}

func TestAllFalse_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	actual := itr.All(x1, func(p Person) bool { return p.Age > 21 })

	assertResult(t, false, actual)
}

func TestMaxEmpty_Enm(t *testing.T) {
	nums := personSlice0()

	x1 := enm.FromSlice(&nums)
	x2 := enm.Select(x1, personAge)
	_, ok := enm.Max(x2)

	assertResult(t, false, ok)
}

func TestMaxEmpty_Itr(t *testing.T) {
	nums := personSlice0()

	x1 := itr.FromSlice(&nums)
	x2 := itr.Select(x1, personAge)
	_, ok := itr.Max(x2)

	assertResult(t, false, ok)
}

func TestMaxOneItem_Enm(t *testing.T) {
	nums := personSlice1()

	x1 := enm.FromSlice(&nums)
	x2 := enm.Select(x1, personAge)
	max, ok := enm.Max(x2)

	assertResult(t, 23, max)
	assertResult(t, true, ok)
}

func TestMaxOneItem_Itr(t *testing.T) {
	nums := personSlice1()

	x1 := itr.FromSlice(&nums)
	x2 := itr.Select(x1, personAge)
	max, ok := itr.Max(x2)

	assertResult(t, 23, max)
	assertResult(t, true, ok)
}

func TestMaxString_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	max, ok := enm.Max(x2)

	assertResult(t, "Zack", max)
	assertResult(t, true, ok)
}

func TestMaxString_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	max, ok := itr.Max(x2)

	assertResult(t, "Zack", max)
	assertResult(t, true, ok)
}

func TestMinEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	_, ok := enm.Min(x2)

	assertResult(t, false, ok)
}

func TestMinEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	_, ok := itr.Min(x2)

	assertResult(t, false, ok)
}

func TestMinOneItem_Enm(t *testing.T) {
	p := personSlice1()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	min, ok := enm.Min(x2)

	assertResult(t, 23, min)
	assertResult(t, true, ok)
}

func TestMinOneItem_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	min, ok := itr.Min(x2)

	assertResult(t, 23, min)
	assertResult(t, true, ok)
}

func TestMinString_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	min, ok := enm.Min(x2)

	assertResult(t, "Abi", min)
	assertResult(t, true, ok)
}

func TestMinString_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	min, ok := itr.Min(x2)

	assertResult(t, "Abi", min)
	assertResult(t, true, ok)
}

func TestAvgEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	_, ok := enm.Avg(x2)

	assertResult(t, false, ok)
}

func TestAvgEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	_, ok := itr.Avg(x2)

	assertResult(t, false, ok)
}

func TestAvgOneItem_Enm(t *testing.T) {
	p := personSlice1()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	avg, ok := enm.Avg(x2)

	assertResult(t, 23.0, avg)
	assertResult(t, true, ok)
}

func TestAvgOneItem_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	avg, ok := itr.Avg(x2)

	assertResult(t, 23.0, avg)
	assertResult(t, true, ok)
}

func TestAvgManyItems_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	avg, ok := enm.Avg(x2)

	assertResult(t, 29.8, avg)
	assertResult(t, true, ok)
}

func TestAvgManyItems_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	avg, ok := itr.Avg(x2)

	assertResult(t, 29.8, avg)
	assertResult(t, true, ok)
}

func TestCountEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	count := enm.Count(x2)

	assertResult(t, 0, count)
}

func TestCountEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	count := itr.Count(x2)

	assertResult(t, 0, count)
}

func TestCountOneItem_Enm(t *testing.T) {
	p := personSlice1()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	count := enm.Count(x2)

	assertResult(t, 1, count)
}

func TestCountOneItem_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	count := itr.Count(x2)

	assertResult(t, 1, count)
}

func TestCount5_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	count := enm.Count(x2)

	assertResult(t, 5, count)
}

func TestCount5_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	count := itr.Count(x2)

	assertResult(t, 5, count)
}

func TestSumEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	count := enm.Count(x2)

	assertResult(t, 0, count)
}

func TestSumEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	count := itr.Count(x2)

	assertResult(t, 0, count)
}

func TestSumOneItem_Enm(t *testing.T) {
	p := personSlice1()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	sum := enm.Sum(x2)

	assertResult(t, 23, sum)
}

func TestSumOneItem_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	sum := itr.Sum(x2)

	assertResult(t, 23, sum)
}

func TestSum5_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	sum := enm.Sum(x2)

	assertResult(t, 149, sum)
}

func TestSum5_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	sum := itr.Sum(x2)

	assertResult(t, 149, sum)
}

func TestSumFiltered_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Where(x1, func(p Person) bool { return p.Age > 30 })
	x3 := enm.Select(x2, personAge)
	sum := enm.Sum(x3)

	assertResult(t, 107, sum)
}

func TestSumFiltered_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Where(x1, func(p Person) bool { return p.Age > 30 })
	x3 := itr.Select(x2, personAge)
	sum := itr.Sum(x3)

	assertResult(t, 107, sum)
}

func TestAccumulateEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	accum := enm.Accumulate(x1, 99, func(prior int, p Person) int { return prior + p.Age })

	assertResult(t, 99, accum)
}

func TestAccumulateEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	accum := itr.Accumulate(x1, 99, func(prior int, p Person) int { return prior + p.Age })

	assertResult(t, 99, accum)
}

func TestAccumulateManyInts_Enm(t *testing.T) {
	nums := intRange(2, 4)

	x1 := enm.FromSlice(&nums)
	accum := enm.Accumulate(x1, -100, func(prior int, item int) int { return prior * item })

	assertResult(t, -2400, accum) // -100 * 2 * 3 * 4
}

func TestAccumulateManyInts_Itr(t *testing.T) {
	nums := intRange(2, 4)

	x1 := itr.FromSlice(&nums)
	accum := itr.Accumulate(x1, -100, func(prior int, item int) int { return prior * item })

	assertResult(t, -2400, accum) // -100 * 2 * 3 * 4
}

func TestAccumulateManyStrings_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Where(x1, func(p Person) bool { return p.Age < 30 })
	accum := enm.Accumulate(x2, "PeopleUnder30:", func(prior string, p Person) string { return prior + p.Name })

	assertResult(t, "PeopleUnder30:JamesAbi", accum) // -100 * 2 * 3 * 4
}

func TestAccumulateManyStrings_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Where(x1, func(p Person) bool { return p.Age < 30 })
	accum := itr.Accumulate(x2, "PeopleUnder30:", func(prior string, p Person) string { return prior + p.Name })

	assertResult(t, "PeopleUnder30:JamesAbi", accum) // -100 * 2 * 3 * 4
}

func TestContainsEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	contains := enm.Contains(x2, 99)

	assertResult(t, false, contains)
}

func TestContainsEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	contains := itr.Contains(x2, 99)

	assertResult(t, false, contains)
}

func TestContainsMany_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	contains999 := enm.Contains(x2, 999)
	contains23 := enm.Contains(x2, 23)

	assertResult(t, false, contains999)
	assertResult(t, true, contains23)
}

func TestContainsMany_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	contains999 := itr.Contains(x2, 999)
	contains23 := itr.Contains(x2, 23)

	assertResult(t, false, contains999)
	assertResult(t, true, contains23)
}

func TestElementAt_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
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

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	element0, ok0 := itr.ElementAt(x2, 0)
	element2, ok2 := itr.ElementAt(x2, 2)
	element4, ok4 := itr.ElementAt(x2, 4)
	element5, ok5 := itr.ElementAt(x2, 5)

	assertResult(t, true, ok0)
	assertResult(t, true, ok2)
	assertResult(t, true, ok4)
	assertResult(t, false, ok5)

	assertResult(t, "James", element0)
	assertResult(t, "Zack", element2)
	assertResult(t, "Rach", element4)
	assertResult(t, "", element5)
}

func TestFirstEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.First(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestFirstEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first, ok := itr.First(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestFirst_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.First(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestFirst_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first, ok := itr.First(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestFirstOrDefaultEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first := enm.FirstOrDefault(x2)

	assertResult(t, "", first)
}

func TestFirstOrDefaultEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first := itr.FirstOrDefault(x2)

	assertResult(t, "", first)
}

func TestFirstOrDefault_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first := enm.FirstOrDefault(x2)

	assertResult(t, "James", first)
}

func TestFirstOrDefault_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first := itr.FirstOrDefault(x2)

	assertResult(t, "James", first)
}

func TestSingleEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingleEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first, ok := itr.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingle_Enm(t *testing.T) {
	p := personSlice1()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Single(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestSingle_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first, ok := itr.Single(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestSingleMany_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingleMany_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first, ok := itr.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingleOrDefaultEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.SingleOrDefault(x2)

	assertResult(t, "", first)
	assertResult(t, true, ok)
}

func TestSingleOrDefaultEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first, ok := itr.SingleOrDefault(x2)

	assertResult(t, "", first)
	assertResult(t, true, ok)
}

func TestSingleOrDefault_Enm(t *testing.T) {
	p := personSlice1()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.SingleOrDefault(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestSingleOrDefault_Itr(t *testing.T) {
	p := personSlice1()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first, ok := itr.SingleOrDefault(x2)

	assertResult(t, "James", first)
	assertResult(t, true, ok)
}

func TestSingleOrDefaultMany_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestSingleOrDefaultMany_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	first, ok := itr.Single(x2)

	assertResult(t, "", first)
	assertResult(t, false, ok)
}

func TestLastEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	last, ok := enm.Last(x2)

	assertResult(t, "", last)
	assertResult(t, false, ok)
}

func TestLastEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	last, ok := itr.Last(x2)

	assertResult(t, "", last)
	assertResult(t, false, ok)
}

func TestLast_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	first, ok := enm.Last(x2)

	assertResult(t, "Rach", first)
	assertResult(t, true, ok)
}

func TestLast_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	last, ok := itr.Last(x2)

	assertResult(t, "Rach", last)
	assertResult(t, true, ok)
}

func TestLastOrDefaultEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	last := enm.LastOrDefault(x2)

	assertResult(t, "", last)
}

func TestLastOrDefaultEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	last := itr.LastOrDefault(x2)

	assertResult(t, "", last)
}

func TestLastOrDefault_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	last := enm.LastOrDefault(x2)

	assertResult(t, "Rach", last)
}

func TestLastOrDefault_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	last := itr.LastOrDefault(x2)

	assertResult(t, "Rach", last)
}

func TestChunkEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	x3 := enm.Chunk(x2, 3)
	chunks := enm.ToSlice(x3)

	assertResult(t, [][]string{}, chunks)
}

func TestChunkEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	x3 := itr.Chunk(x2, 3)
	chunks := itr.ToSlice(x3)

	assertResult(t, [][]string{}, chunks)
}

func TestChunk_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personName)
	x3 := enm.Chunk(x2, 3)
	chunks := enm.ToSlice(x3)

	expected := [][]string{
		{"James", "Lucy", "Zack"},
		{"Abi", "Rach"},
	}
	assertResult(t, expected, chunks)
}

func TestChunk_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personName)
	x3 := itr.Chunk(x2, 3)
	chunks := itr.ToSlice(x3)

	expected := [][]string{
		{"James", "Lucy", "Zack"},
		{"Abi", "Rach"},
	}
	assertResult(t, expected, chunks)
}

func TestDistinctEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	x3 := enm.Distinct(x2)
	distinct := enm.ToSlice(x3)

	assertResult(t, []int{}, distinct)
}

func TestDistinctEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	x3 := itr.Distinct(x2)
	distinct := itr.ToSlice(x3)

	assertResult(t, []int{}, distinct)
}

func TestDistinct_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	x2 := enm.Select(x1, personAge)
	x3 := enm.Distinct(x2)
	distinct := enm.ToSlice(x3)

	assertResult(t, []int{23, 33, 41, 19}, distinct)
}

func TestDistinct_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	x2 := itr.Select(x1, personAge)
	x3 := itr.Distinct(x2)
	distinct := itr.ToSlice(x3)

	assertResult(t, []int{23, 33, 41, 19}, distinct)
}

func TestToMap_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	mapped, ok := enm.ToMap(x1,
		func(p Person) string { return strings.ToUpper(p.Name) },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, true)
	assertResult(t, mapped, map[string]int{"JAMES": 230, "LUCY": 330, "ZACK": 410, "ABI": 190, "RACH": 330})
}

func TestToMap_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	mapped, ok := itr.ToMap(x1,
		func(p Person) string { return strings.ToUpper(p.Name) },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, true)
	assertResult(t, mapped, map[string]int{"JAMES": 230, "LUCY": 330, "ZACK": 410, "ABI": 190, "RACH": 330})
}

func TestToMapEmpty_Enm(t *testing.T) {
	p := personSlice0()

	x1 := enm.FromSlice(&p)
	mapped, ok := enm.ToMap(x1,
		func(p Person) string { return strings.ToUpper(p.Name) },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, true)
	assertResult(t, len(mapped), 0)
}

func TestToMapEmpty_Itr(t *testing.T) {
	p := personSlice0()

	x1 := itr.FromSlice(&p)
	mapped, ok := itr.ToMap(x1,
		func(p Person) string { return strings.ToUpper(p.Name) },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, true)
	assertResult(t, len(mapped), 0)
}

func TestToMapDuplicateKey_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.FromSlice(&p)
	mapped, ok := enm.ToMap(x1,
		func(p Person) string { return "DUPLICATE" },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, false)
	assertResult(t, mapped, nil)
}

func TestToMapDuplicateKey_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.FromSlice(&p)
	mapped, ok := itr.ToMap(x1,
		func(p Person) string { return "DUPLICATE" },
		func(p Person) int { return p.Age * 10 })

	assertResult(t, ok, false)
	assertResult(t, mapped, nil)
}

func TestPointersFromSlice_Enm(t *testing.T) {
	p := personSlice5()

	x1 := enm.PointersFromSlice(&p)
	x2 := enm.Select(x1, func(ptr *Person) int { return (*ptr).Age })
	maxAge, _ := enm.Max(x2)

	assertResult(t, maxAge, 41)
}

func TestPointersFromSlice_Itr(t *testing.T) {
	p := personSlice5()

	x1 := itr.PointersFromSlice(&p)
	x2 := itr.Select(x1, func(ptr *Person) int { return (*ptr).Age })
	maxAge, _ := itr.Max(x2)

	assertResult(t, maxAge, 41)
}

func TestFromMap_Enm(t *testing.T) {
	p := personMap5()

	x1 := enm.FromMap(&p)
	x2 := enm.Where(x1, func(kvp cmn.KeyValuePair[string, Person]) bool { return kvp.Key == "Abi" })
	x3 := enm.Select(x2, func(kvp cmn.KeyValuePair[string, Person]) int { return kvp.Value.Age })
	abiAge, _ := enm.Single(x3)

	assertResult(t, abiAge, 19)
}

func TestFromMap_Itr(t *testing.T) {
	p := personMap5()

	x1 := itr.FromMap(&p)
	x2 := itr.Where(x1, func(kvp cmn.KeyValuePair[string, Person]) bool { return kvp.Key == "Abi" })
	x3 := itr.Select(x2, func(kvp cmn.KeyValuePair[string, Person]) int { return kvp.Value.Age })
	abiAge, _ := itr.Single(x3)

	assertResult(t, abiAge, 19)
}

func TestFromMapEmpty_Enm(t *testing.T) {
	p := personMap0()

	x1 := enm.FromMap(&p)
	slice := enm.ToSlice(x1)

	assertResult(t, len(slice), 0)
}

func TestFromMapEmpty_Itr(t *testing.T) {
	p := personMap0()

	x1 := itr.FromMap(&p)
	slice := itr.ToSlice(x1)

	assertResult(t, len(slice), 0)
}

func TestPointersFromMap_Enm(t *testing.T) {
	//TODO - intermittent failures
	p := personMap5()

	x1 := enm.PointersFromMap(&p)
	x2 := enm.Where(x1, func(kvp cmn.KeyValuePair[string, *Person]) bool { return kvp.Key == "Abi" })
	x3 := enm.Select(x2, func(kvp cmn.KeyValuePair[string, *Person]) int { return (*kvp.Value).Age })
	abiAge, _ := enm.Single(x3)

	assertResult(t, abiAge, 19)
}

func TestPointersFromMap_Itr(t *testing.T) {
	//TODO - intermittent failures
	p := personMap5()

	x1 := itr.PointersFromMap(&p)
	x2 := itr.Where(x1, func(kvp cmn.KeyValuePair[string, *Person]) bool { return kvp.Key == "Abi" })
	x3 := itr.Select(x2, func(kvp cmn.KeyValuePair[string, *Person]) int { return (*kvp.Value).Age })
	abiAge, _ := itr.Single(x3)

	assertResult(t, abiAge, 19)
}

func TestPointersFromMapEmpty_Enm(t *testing.T) {
	p := personMap0()

	x1 := enm.PointersFromMap(&p)
	slice := enm.ToSlice(x1)

	assertResult(t, len(slice), 0)
}

func TestPointersFromMapEmpty_Itr(t *testing.T) {
	p := personMap0()

	x1 := itr.PointersFromMap(&p)
	slice := itr.ToSlice(x1)

	assertResult(t, len(slice), 0)
}
