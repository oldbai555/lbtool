package pie

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// Abs is a function which returns the absolute value of all the
// elements in the slice.
func (ss Uint32s) Abs() Uint32s {
	result := make(Uint32s, len(ss))
	for i, val := range ss {
		if val < 0 {
			result[i] = -val
		} else {
			result[i] = val
		}
	}
	return result
}

// All will return true if all callbacks return true. It follows the same logic
// as the all() function in Python.
//
// If the list is empty then true is always returned.
func (ss Uint32s) All(fn func(value uint32) bool) bool {
	for _, value := range ss {
		if !fn(value) {
			return false
		}
	}

	return true
}

// Any will return true if any callbacks return true. It follows the same logic
// as the any() function in Python.
//
// If the list is empty then false is always returned.
func (ss Uint32s) Any(fn func(value uint32) bool) bool {
	for _, value := range ss {
		if fn(value) {
			return true
		}
	}

	return false
}

// Append will return a new slice with the elements appended to the end.
//
// It is acceptable to provide zero arguments.
func (ss Uint32s) Append(elements ...uint32) Uint32s {
	// Copy ss, to make sure no memory is overlapping between input and
	// output. See issue #97.
	result := append(Uint32s{}, ss...)

	result = append(result, elements...)
	return result
}

// AreSorted will return true if the slice is already sorted. It is a wrapper
// for sort.Uint32sAreSorted.
func (ss Uint32s) AreSorted() bool {
	return sort.SliceIsSorted(ss, func(i, j int) bool {
		return ss[i] < ss[j]
	})
}

// AreUnique will return true if the slice contains elements that are all
// different (unique) from each other.
func (ss Uint32s) AreUnique() bool {
	return ss.Unique().Len() == ss.Len()
}

// Average is the average of all of the elements, or zero if there are no
// elements.
func (ss Uint32s) Average() float64 {
	if l := uint32(len(ss)); l > 0 {
		return float64(ss.Sum()) / float64(l)
	}

	return 0
}

// Bottom will return n elements from bottom
//
// that means that elements is taken from the end of the slice
// for this [1,2,3] slice with n == 2 will be returned [3,2]
// if the slice has less elements then n that'll return all elements
// if n < 0 it'll return empty slice.
func (ss Uint32s) Bottom(n int) (top Uint32s) {
	var lastIndex = len(ss) - 1
	for i := lastIndex; i > -1 && n > 0; i-- {
		top = append(top, ss[i])
		n--
	}

	return
}

// Contains returns true if the element exists in the slice.
//
// When using slices of pointers it will only compare by address, not value.
func (ss Uint32s) Contains(lookingFor uint32) bool {
	for _, s := range ss {
		if lookingFor == s {
			return true
		}
	}

	return false
}

// Diff returns the elements that needs to be added or removed from the first
// slice to have the same elements in the second slice.
//
// The order of elements is not taken into consideration, so the slices are
// treated sets that allow duplicate items.
//
// The added and removed returned may be blank respectively, or contain upto as
// many elements that exists in the largest slice.
func (ss Uint32s) Diff(against Uint32s) (added, removed Uint32s) {
	// This is probably not the best way to do it. We do an O(n^2) between the
	// slices to see which items are missing in each direction.

	diffOneWay := func(ss1, ss2raw Uint32s) (result Uint32s) {
		ss2 := make(Uint32s, len(ss2raw))
		copy(ss2, ss2raw)

		for _, s := range ss1 {
			found := false

			for i, element := range ss2 {
				if s == element {
					ss2 = append(ss2[:i], ss2[i+1:]...)
					found = true
					break
				}
			}

			if !found {
				result = append(result, s)
			}
		}

		return
	}

	removed = diffOneWay(ss, against)
	added = diffOneWay(against, ss)

	return
}

// DropTop will return the rest slice after dropping the top n elements
// if the slice has less elements then n that'll return empty slice
// if n < 0 it'll return empty slice.
func (ss Uint32s) DropTop(n int) (drop Uint32s) {
	if n < 0 || n >= len(ss) {
		return
	}

	// Copy ss, to make sure no memory is overlapping between input and
	// output. See issue #145.
	drop = make([]uint32, len(ss)-n)
	copy(drop, ss[n:])

	return
}

// Drop items from the slice while f(item) is true.
// Afterwards, return every element until the slice is empty. It follows the same logic as the dropwhile() function from itertools in Python.
func (ss Uint32s) DropWhile(f func(s uint32) bool) (ss2 Uint32s) {
	ss2 = make([]uint32, len(ss))
	copy(ss2, ss)
	for i, value := range ss2 {
		if !f(value) {
			return ss2[i:]
		}
	}
	return Uint32s{}
}

// Each is more condensed version of Transform that allows an action to happen
// on each elements and pass the original slice on.
//
//	cars.Each(func (car *Car) {
//	    fmt.Printf("Car color is: %s\n", car.Color)
//	})
//
// Pie will not ensure immutability on items passed in so they can be
// manipulated, if you choose to do it this way, for example:
//
//	// Set all car colors to Red.
//	cars.Each(func (car *Car) {
//	    car.Color = "Red"
//	})
func (ss Uint32s) Each(fn func(uint32)) Uint32s {
	for _, s := range ss {
		fn(s)
	}

	return ss
}

// Equals compare elements from the start to the end,
//
// if they are the same is considered the slices are equal if all elements are the same is considered the slices are equal
// if each slice == nil is considered that they're equal
//
// if element realizes Equals interface it uses that method, in other way uses default compare
func (ss Uint32s) Equals(rhs Uint32s) bool {
	if len(ss) != len(rhs) {
		return false
	}

	for i := range ss {
		if !(ss[i] == rhs[i]) {
			return false
		}
	}

	return true
}

// Extend will return a new slice with the slices of elements appended to the
// end.
//
// It is acceptable to provide zero arguments.
func (ss Uint32s) Extend(slices ...Uint32s) (ss2 Uint32s) {
	ss2 = ss

	for _, slice := range slices {
		ss2 = ss2.Append(slice...)
	}

	return ss2
}

// Filter will return a new slice containing only the elements that return
// true from the condition. The returned slice may contain zero elements (nil).
//
// FilterNot works in the opposite way of Filter.
func (ss Uint32s) Filter(condition func(uint32) bool) (ss2 Uint32s) {
	for _, s := range ss {
		if condition(s) {
			ss2 = append(ss2, s)
		}
	}
	return
}

// FilterNot works the same as Filter, with a negated condition. That is, it will
// return a new slice only containing the elements that returned false from the
// condition. The returned slice may contain zero elements (nil).
func (ss Uint32s) FilterNot(condition func(uint32) bool) (ss2 Uint32s) {
	for _, s := range ss {
		if !condition(s) {
			ss2 = append(ss2, s)
		}
	}

	return
}

// FindFirstUsing will return the index of the first element when the callback returns true or -1 if no element is found.
// It follows the same logic as the findIndex() function in Javascript.
//
// If the list is empty then -1 is always returned.
func (ss Uint32s) FindFirstUsing(fn func(value uint32) bool) int {
	for idx, value := range ss {
		if fn(value) {
			return idx
		}
	}

	return -1
}

// First returns the first element, or zero. Also see FirstOr().
func (ss Uint32s) First() uint32 {
	return ss.FirstOr(0)
}

// FirstOr returns the first element or a default value if there are no
// elements.
func (ss Uint32s) FirstOr(defaultValue uint32) uint32 {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[0]
}

// Float64s transforms each element to a float64.
func (ss Uint32s) Float64s() Float64s {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(Float64s, l)
	for i := 0; i < l; i++ {
		mightBeString := ss[i]
		result[i], _ = strconv.ParseFloat(fmt.Sprintf("%v", mightBeString), 64)
	}

	return result
}

// Group returns a map of the value with an individual count.
func (ss Uint32s) Group() map[uint32]int {
	group := map[uint32]int{}
	for _, n := range ss {
		group[n]++
	}
	return group
}

// Intersect returns items that exist in all lists.
//
// It returns slice without any duplicates.
// If zero slice arguments are provided, then nil is returned.
func (ss Uint32s) Intersect(slices ...Uint32s) (ss2 Uint32s) {
	if slices == nil {
		return nil
	}

	var uniqs = make([]map[uint32]struct{}, len(slices))
	for i := 0; i < len(slices); i++ {
		m := make(map[uint32]struct{})
		for _, el := range slices[i] {
			m[el] = struct{}{}
		}
		uniqs[i] = m
	}

	var containsInAll = false
	for _, el := range ss.Unique() {
		for _, u := range uniqs {
			if _, exists := u[el]; !exists {
				containsInAll = false
				break
			}
			containsInAll = true
		}
		if containsInAll {
			ss2 = append(ss2, el)
		}
	}

	return
}

// Insert a value at an index
func (ss Uint32s) Insert(index int, values ...uint32) Uint32s {
	if index >= ss.Len() {
		return Uint32s.Extend(ss, Uint32s(values))
	}

	return Uint32s.Extend(ss[:index], Uint32s(values), ss[index:])
}

// Ints transforms each element to an integer.
func (ss Uint32s) Ints() Ints {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(Ints, l)
	for i := 0; i < l; i++ {
		mightBeString := ss[i]
		f, _ := strconv.ParseFloat(fmt.Sprintf("%v", mightBeString), 64)
		result[i] = int(f)
	}

	return result
}

// Join returns a string from joining each of the elements.
func (ss Uint32s) Join(glue string) (s string) {
	var slice interface{} = []uint32(ss)

	if y, ok := slice.([]string); ok {
		// The stdlib is efficient for type []string
		return strings.Join(y, glue)
	} else {
		// General case
		parts := make([]string, len(ss))
		for i, element := range ss {
			mightBeString := element
			parts[i] = fmt.Sprintf("%v", mightBeString)
		}
		return strings.Join(parts, glue)
	}
}

// JSONBytes returns the JSON encoded array as bytes.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array.
func (ss Uint32s) JSONBytes() []byte {
	if ss == nil {
		return []byte("[]")
	}

	// An error should not be possible.
	data, _ := json.Marshal(ss)

	return data
}

// JSONBytesIndent returns the JSON encoded array as bytes with indent applied.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array. See
// json.MarshalIndent for details.
func (ss Uint32s) JSONBytesIndent(prefix, indent string) []byte {
	if ss == nil {
		return []byte("[]")
	}

	// An error should not be possible.
	data, _ := json.MarshalIndent(ss, prefix, indent)

	return data
}

// JSONString returns the JSON encoded array as a string.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array.
func (ss Uint32s) JSONString() string {
	if ss == nil {
		return "[]"
	}

	// An error should not be possible.
	data, _ := json.Marshal(ss)

	return string(data)
}

// JSONStringIndent returns the JSON encoded array as a string with indent applied.
//
// One important thing to note is that it will treat a nil slice as an empty
// slice to ensure that the JSON value return is always an array. See
// json.MarshalIndent for details.
func (ss Uint32s) JSONStringIndent(prefix, indent string) string {
	if ss == nil {
		return "[]"
	}

	// An error should not be possible.
	data, _ := json.MarshalIndent(ss, prefix, indent)

	return string(data)
}

// Last returns the last element, or zero. Also see LastOr().
func (ss Uint32s) Last() uint32 {
	return ss.LastOr(0)
}

// LastOr returns the last element or a default value if there are no elements.
func (ss Uint32s) LastOr(defaultValue uint32) uint32 {
	if len(ss) == 0 {
		return defaultValue
	}

	return ss[len(ss)-1]
}

// Len returns the number of elements.
func (ss Uint32s) Len() int {
	return len(ss)
}

// Map will return a new slice where each element has been mapped (transformed).
// The number of elements returned will always be the same as the input.
//
// Be careful when using this with slices of pointers. If you modify the input
// value it will affect the original slice. Be sure to return a new allocated
// object or deep copy the existing one.
func (ss Uint32s) Map(fn func(uint32) uint32) (ss2 Uint32s) {
	if ss == nil {
		return nil
	}

	ss2 = make([]uint32, len(ss))
	for i, s := range ss {
		ss2[i] = fn(s)
	}

	return
}

// Max is the maximum value, or zero.
func (ss Uint32s) Max() (max uint32) {
	if len(ss) == 0 {
		return
	}

	max = ss[0]
	for _, s := range ss {
		if s > max {
			max = s
		}
	}

	return
}

// Median returns the value separating the higher half from the lower half of a
// data sample.
//
// Zero is returned if there are no elements in the slice.
//
// If the number of elements is even, then the uint32 mean of the two "median values"
// is returned.
func (ss Uint32s) Median() uint32 {
	n := len(ss)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return ss[0]
	}

	// This implementation aims at linear time O(n) on average.
	// It uses the same idea as QuickSort, but makes only 1 recursive
	// call instead of 2. See also Quickselect.

	work := make(Uint32s, len(ss))
	copy(work, ss)

	limit1, limit2 := n/2, n/2+1
	if n%2 == 0 {
		limit1, limit2 = n/2-1, n/2+1
	}

	var rec func(a, b int)
	rec = func(a, b int) {
		if b-a <= 1 {
			return
		}
		ipivot := (a + b) / 2
		pivot := work[ipivot]
		work[a], work[ipivot] = work[ipivot], work[a]
		j := a
		k := b
		for j+1 < k {
			if work[j+1] < pivot {
				work[j+1], work[j] = work[j], work[j+1]
				j++
			} else {
				work[j+1], work[k-1] = work[k-1], work[j+1]
				k--
			}
		}
		// 1 or 0 recursive calls
		if j > limit1 {
			rec(a, j)
		}
		if j+1 < limit2 {
			rec(j+1, b)
		}
	}

	rec(0, len(work))

	if n%2 == 1 {
		return work[n/2]
	} else {
		return (work[n/2-1] + work[n/2]) / 2
	}
}

// Min is the minimum value, or zero.
func (ss Uint32s) Min() (min uint32) {
	if len(ss) == 0 {
		return
	}

	min = ss[0]
	for _, s := range ss {
		if s < min {
			min = s
		}
	}

	return
}

// Mode returns a new slice containing the most frequently occuring values.
//
// The number of items returned may be the same as the input or less. It will
// never return zero items unless the input slice has zero items.
func (ss Uint32s) Mode() Uint32s {
	if len(ss) == 0 {
		return nil
	}
	values := make(map[uint32]int)
	for _, s := range ss {
		values[s]++
	}

	var maxFrequency int
	for _, v := range values {
		if v > maxFrequency {
			maxFrequency = v
		}
	}

	var maxValues Uint32s
	for k, v := range values {
		if v == maxFrequency {
			maxValues = append(maxValues, k)
		}
	}

	return maxValues
}

// Pop the first element of the slice
//
// Usage Example:
//
//	type knownGreetings []string
//	greetings := knownGreetings{"ciao", "hello", "hola"}
//	for greeting := greetings.Pop(); greeting != nil; greeting = greetings.Pop() {
//	    fmt.Println(*greeting)
//	}
func (ss *Uint32s) Pop() (popped *uint32) {

	if len(*ss) == 0 {
		return
	}

	popped = &(*ss)[0]
	*ss = (*ss)[1:]
	return
}

// Product is the product of all of the elements.
func (ss Uint32s) Product() (product uint32) {
	if len(ss) == 0 {
		return
	}
	product = ss[0]
	for _, s := range ss[1:] {
		product *= s
	}

	return
}

// Random returns a random element by your rand.Source, or zero
func (ss Uint32s) Random(source rand.Source) uint32 {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 1 {
		return 0
	}
	if n < 2 {
		return ss[0]
	}
	rnd := rand.New(source)
	i := rnd.Intn(n)
	return ss[i]
}

// Reduce continually applies the provided function
// over the slice. Reducing the elements to a single value.
//
// Returns a zero value of uint32 if there are no elements in the slice. It will panic if the reducer is nil and the slice has more than one element (required to invoke reduce).
// Otherwise returns result of applying reducer from left to right.
func (ss Uint32s) Reduce(reducer func(uint32, uint32) uint32) (el uint32) {
	if len(ss) == 0 {
		return
	}
	el = ss[0]
	for _, s := range ss[1:] {
		el = reducer(el, s)
	}
	return
}

// Reverse returns a new copy of the slice with the elements ordered in reverse.
// This is useful when combined with Sort to get a descending sort order:
//
//	ss.Sort().Reverse()
func (ss Uint32s) Reverse() Uint32s {
	// Avoid the allocation. If there is one element or less it is already
	// reversed.
	if len(ss) < 2 {
		return ss
	}

	sorted := make([]uint32, len(ss))
	for i := 0; i < len(ss); i++ {
		sorted[i] = ss[len(ss)-i-1]
	}

	return sorted
}

// Send sends elements to channel
// in normal act it sends all elements but if func canceled it can be less
//
// it locks execution of gorutine
// it doesn't close channel after work
// returns sended elements if len(this) != len(old) considered func was canceled
func (ss Uint32s) Send(ctx context.Context, ch chan<- uint32) Uint32s {
	for i, s := range ss {
		select {
		case <-ctx.Done():
			return ss[:i]
		default:
			ch <- s
		}
	}

	return ss
}

// Sequence generates all numbers in range or returns nil if params invalid
//
// There are 3 variations to generate:
//  1. [0, n).
//  2. [min, max).
//  3. [min, max) with step.
//
// if len(params) == 1 considered that will be returned slice between 0 and n,
// where n is the first param, [0, n).
// if len(params) == 2 considered that will be returned slice between min and max,
// where min is the first param, max is the second, [min, max).
// if len(params) > 2 considered that will be returned slice between min and max with step,
// where min is the first param, max is the second, step is the third one, [min, max) with step,
// others params will be ignored
func (ss Uint32s) Sequence(params ...int) Uint32s {
	var creator = func(i int) uint32 {
		return uint32(i)
	}

	return ss.SequenceUsing(creator, params...)
}

// SequenceUsing generates slice in range using creator function
//
// There are 3 variations to generate:
//  1. [0, n).
//  2. [min, max).
//  3. [min, max) with step.
//
// if len(params) == 1 considered that will be returned slice between 0 and n,
// where n is the first param, [0, n).
// if len(params) == 2 considered that will be returned slice between min and max,
// where min is the first param, max is the second, [min, max).
// if len(params) > 2 considered that will be returned slice between min and max with step,
// where min is the first param, max is the second, step is the third one, [min, max) with step,
// others params will be ignored
func (ss Uint32s) SequenceUsing(creator func(int) uint32, params ...int) Uint32s {
	var seq = func(min, max, step int) (seq Uint32s) {
		lenght := int(Round(float64(max-min) / float64(step)))
		if lenght < 1 {
			return
		}

		seq = make(Uint32s, lenght)
		for i := 0; i < lenght; min += step {
			seq[i] = creator(min)
			i++
		}

		return seq
	}

	if len(params) > 2 {
		return seq(params[0], params[1], params[2])
	} else if len(params) == 2 {
		return seq(params[0], params[1], 1)
	} else if len(params) == 1 {
		return seq(0, params[0], 1)
	} else {
		return nil
	}
}

// Shift will return two values: the shifted value and the rest slice.
func (ss Uint32s) Shift() (uint32, Uint32s) {
	return ss.First(), ss.DropTop(1)
}

// Shuffle returns shuffled slice by your rand.Source
func (ss Uint32s) Shuffle(source rand.Source) Uint32s {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 2 {
		return ss
	}

	// go 1.10+ provides rnd.Shuffle. However, to support older versions we copy
	// the algorithm directly from the go source: src/math/rand/rand.go below,
	// with some adjustments:
	shuffled := make([]uint32, n)
	copy(shuffled, ss)

	rnd := rand.New(source)

	Shuffle(rnd, n, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

// Sort works similar to sort.Uint32s(). However, unlike sort.Uint32s the
// slice returned will be reallocated as to not modify the input slice.
//
// See Reverse() and AreSorted().
func (ss Uint32s) Sort() Uint32s {
	// Avoid the allocation. If there is one element or less it is already
	// sorted.
	if len(ss) < 2 {
		return ss
	}

	sorted := make(Uint32s, len(ss))
	copy(sorted, ss)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	return sorted
}

// Stddev is the standard deviation
func (ss Uint32s) Stddev() float64 {
	if len(ss) == 0 {
		return 0.0
	}

	avg := ss.Average()

	var sd float64
	for i := range ss {
		sd += math.Pow(float64(ss[i])-avg, 2)
	}
	sd = math.Sqrt(sd / float64(len(ss)))

	return sd
}

// Strings transforms each element to a string.
//
// If the element type implements fmt.Stringer it will be used. Otherwise it
// will fallback to the result of:
//
//	fmt.Sprintf("%v")
func (ss Uint32s) Strings() Strings {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(Strings, l)
	for i := 0; i < l; i++ {
		mightBeString := ss[i]
		result[i] = fmt.Sprintf("%v", mightBeString)
	}

	return result
}

// SubSlice will return the subSlice from start to end(excluded)
//
// Condition 1: If start < 0 or end < 0, nil is returned.
// Condition 2: If start >= end, nil is returned.
// Condition 3: Return all elements that exist in the range provided,
// if start or end is out of bounds, zero items will be placed.
func (ss Uint32s) SubSlice(start int, end int) (subSlice Uint32s) {
	if start < 0 || end < 0 {
		return
	}

	if start >= end {
		return
	}

	length := ss.Len()
	if start < length {
		if end <= length {
			subSlice = ss[start:end]
		} else {
			zeroArray := make([]uint32, end-length)
			subSlice = ss[start:length].Append(zeroArray[:]...)
		}
	} else {
		zeroArray := make([]uint32, end-start)
		subSlice = zeroArray[:]
	}

	return
}

// Sum is the sum of all of the elements.
func (ss Uint32s) Sum() (sum uint32) {
	for _, s := range ss {
		sum += s
	}

	return
}

// Top will return n elements from head of the slice
// if the slice has less elements then n that'll return all elements
// if n < 0 it'll return empty slice.
func (ss Uint32s) Top(n int) (top Uint32s) {
	for i := 0; i < len(ss) && n > 0; i++ {
		top = append(top, ss[i])
		n--
	}

	return
}

// StringsUsing transforms each element to a string.
func (ss Uint32s) StringsUsing(transform func(uint32) string) Strings {
	l := len(ss)

	// Avoid the allocation.
	if l == 0 {
		return nil
	}

	result := make(Strings, l)
	for i := 0; i < l; i++ {
		result[i] = transform(ss[i])
	}

	return result
}

// Unique returns a new slice with all of the unique values.
//
// The items will be returned in a randomized order, even with the same input.
//
// The number of items returned may be the same as the input or less. It will
// never return zero items unless then input slice has zero items.
//
// A slice with zero elements is considered to be unique.
//
// See AreUnique().
func (ss Uint32s) Unique() Uint32s {
	// Avoid the allocation. If there is one element or less it is already
	// unique.
	if len(ss) < 2 {
		return ss
	}

	values := map[uint32]struct{}{}

	for _, value := range ss {
		values[value] = struct{}{}
	}

	var uniqueValues Uint32s
	for value := range values {
		uniqueValues = append(uniqueValues, value)
	}

	return uniqueValues
}

// Unshift adds one or more elements to the beginning of the slice
// and returns the new slice.
func (ss Uint32s) Unshift(elements ...uint32) (unshift Uint32s) {
	unshift = append(Uint32s{}, elements...)
	unshift = append(unshift, ss...)

	return
}
