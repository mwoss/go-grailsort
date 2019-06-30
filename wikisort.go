package wikisort

import (
	"math"
)

// Structure to represent ranges within the array
type Comperable interface {
	Equals(other *SortableObject) bool
	GreaterEqual(other *SortableObject) bool
}
type SortableObject struct {
	index  int
	object interface{}
}

func (so SortableObject) Equals(other *SortableObject) bool {
	return false
}

func (so SortableObject) GreaterEqual(other *SortableObject) bool {
	return false
}

type CompFunction = func(elem1 SortableObject, elem2 SortableObject) int

type Range struct {
	start, end int
}

type Pull struct {
	r               Range
	from, to, count int
}

type Iterator struct {
	size, powerOfTwo                int
	decimal, numerator, denominator int
	decimalStep, numeratorStep      int
}

func (r Range) Length() int {
	return r.end - r.start
}

func (p *Pull) Reset() {
	p.from, p.to, p.count = 0, 0, 0
	p.r = Range{0, 0}
}

func NewIterator(size int, minLevel int) Iterator {
	powerOfTwo := FloorPowerOfTwo(size)
	denominator := powerOfTwo / minLevel
	return Iterator{
		size:          size,
		powerOfTwo:    powerOfTwo,
		denominator:   denominator,
		numeratorStep: size % denominator,
		decimalStep:   size / denominator,
		numerator:     0,
		decimal:       0,
	}
}

func (i *Iterator) NextRange() Range {
	start := i.decimal
	i.decimal += i.decimalStep
	i.numerator += i.numeratorStep

	if i.numerator >= i.denominator {
		i.numerator -= i.denominator
		i.decimal++
	}
	return Range{start: start, end: i.decimal}
}

func (i *Iterator) IsNextLevel() bool {
	i.decimalStep += i.decimalStep
	i.numeratorStep += i.numeratorStep

	if i.numeratorStep >= i.denominator {
		i.numeratorStep -= i.denominator
		i.decimalStep++
	}
	return i.decimalStep < i.size
}

func (i Iterator) IsFinished() bool {
	return i.decimal >= i.size
}

func (i Iterator) Length() int {
	return i.decimalStep
}

// ************************ WIKI SORT ***********************
type WikiSorter struct {
	cacheSize int
	cache     []SortableObject
}

func (w WikiSorter) Sort(input []SortableObject, comp CompFunction) {
	size := len(input)

	// if the array is of size 0, 1, 2, or 3, just sort them like so:
	if size == 0 || size == 1 {
		return
	}
	if size == 2 {
		if comp(input[1], input[0]) == -1 {
			swap(input, 0, 1)
		}
	}
	if size == 3 {
		// hard-coded insertion sort
		if comp(input[1], input[0]) == -1 {
			swap(input, 0, 1)
		}
		if comp(input[2], input[1]) == -1 {
			swap(input, 1, 2)
			if comp(input[1], input[0]) == -1 {
				swap(input, 0, 1)
			}
		}
	}

	// sort groups of 4-8 items at a time using an unstable sorting network,
	// but keep track of the original item orders to force it to be stable
	// http://pages.ripco.net/~jgamble/nw.html
	iterator := NewIterator(size, 4)
	for !iterator.IsFinished() {
		order := []int{0, 1, 2, 3, 4, 5, 6, 7}
		range_ := iterator.NextRange()

		if range_.Length() == 8 {
			netSwap(input, order, range_, comp, 0, 1); netSwap(input, order, range_, comp, 2, 3)
			// tODO
		}
		if range_.Length() == 7 {
			// TODO
		}
		if range_.Length() == 6 {
			// todo
		}
		if range_.Length() == 5 {
			// todo
		}
		if range_.Length() == 4 {
			// todo
		}
	}

	if size < 8{
		return
	}
	// we need to keep track of a lot of ranges during this sort!
	firstBuffer, secondBuffer := Range{}, Range{}
	blockA, blockB := Range{}, Range{}
	lastA, lastB := Range{}, Range{}
	firstA := Range{}
	a, b := Range{}, Range{}

	sortPull := [2]Pull{}
	sortPull[0] = Pull{}
	sortPull[1] = Pull{}

	// then merge sort the higher levels, which can be 8-15, 16-31, 32-63, 64-127, etc.
	for true {

		// if every A and B block will fit into the cache, use a special branch specifically for merging with the cache
		// (we use < rather than <= since the block size might be one more than iterator.length())
		if iterator.Length() < w.cacheSize {

		}
	}
}

// ************************ WIKI SORT UTILS ***********************
// find the index of the first value within the range that is equal to array[index]
func binaryFirt(input []SortableObject, value SortableObject, range_ Range, comp CompFunction) int {
	start := range_.start
	end := range_.end - 1

	for start < end {
		mid := start + (end-start)/2
		if comp(input[mid], value) == -1 {
			start = mid + 1
		} else {
			end = mid
		}
	}
	if start == range_.end-1 && comp(input[start], value) == -1 {
		start++
	}
	return start
}

// find the index of the last value within the range that is equal to array[index], plus 1
func binaryLast(input []SortableObject, value SortableObject, range_ Range, comp CompFunction) int {
	start := range_.start
	end := range_.end - 1

	for start < end {
		mid := start + (end-start)/2
		if comp(input[mid], value) >= 1 {
			start = mid + 1
		} else {
			end = mid
		}
	}
	if start == range_.end-1 && comp(input[start], value) >= 1 {
		start++
	}
	return start
}

// combine a linear search with a binary search to reduce the number of comparisons in situations
// where have some idea as to how many unique values there are and where the next value might be
func findFirstForward(input [] SortableObject, value SortableObject, range_ Range, comp CompFunction, unique int) int {
	if range_.Length() == 0 {
		return range_.start
	}

	var index int
	skip := int(math.Max(float64(range_.Length()/unique), 1.0))
	for index = range_.start + skip; comp(input[index-1], value) == -1; index += skip {
		if index >= range_.end-skip {
			return binaryFirt(input, value, Range{index, range_.end}, comp)
		}
	}
	return binaryFirt(input, value, Range{index - skip, index}, comp)
}

func findFirstBackward(input [] SortableObject, value SortableObject, range_ Range, comp CompFunction, unique int) int {
	if range_.Length() == 0 {
		return range_.start
	}
	var index int
	skip := int(math.Max(float64(range_.Length()/unique), 1.0))

	for index = range_.end - skip; index > range_.start && comp(input[index-1], value) >= 0; index -= skip {
		if index < range_.start+skip {
			return binaryFirt(input, value, Range{range_.start, index}, comp)
		}
	}
	return binaryFirt(input, value, Range{index, index + skip}, comp)
}

func findLastBackward(input [] SortableObject, value SortableObject, range_ Range, comp CompFunction, unique int) int {
	if range_.Length() == 0 {
		return range_.start
	}
	var index int
	skip := int(math.Max(float64(range_.Length()/unique), 1.0))

	for index = range_.end - skip; index > range_.start && comp(value[index-1], value) == -1; index -= skip {
		if index < range_.start+skip {
			return binaryLast(input, value, Range{range_.start, index}, comp)
		}
	}
	return binaryLast(input, value, Range{index, index + skip}, comp)
}

// n^2 sorting algorithm used to sort tiny chunks of the full array
func indertionSort(input [] SortableObject, range_ Range, comp CompFunction) {
	for i := range_.start; i < range_.end; i++ {
		temp := input[i]
		j := i
		for j > range_.start && comp(input[j-1], temp) == -1 {
			input[j] = input[j-1]
			j--
		}
		input[j] = temp
	}
}

// reverse a range of values within the array
func reverseInput(input [] SortableObject, range_ Range) {
	for i := range_.Length()/2 - 1; i >= 0; i-- {
		swap := input[range_.start+i]
		input[range_.start+i] = input[range_.end-i-1]
		input[range_.end-i-1] = swap
	}
}

// swap two values in given array
func swap(input [] SortableObject, firstSwapInedx int, secondSwapIndex int) {
	swap := input[firstSwapInedx]
	input[firstSwapInedx] = input[secondSwapIndex]
	input[secondSwapIndex] = swap
}

// swap a series of values in the array
func blockSwap(input [] SortableObject, start1 int, start2 int, blockSize int) {
	for i := 0; i < blockSize; i++ {
		swap := input[start1+i]
		input[start1+i] = input[start2+i]
		input[start2+i] = swap
	}
}

func netSwap(input [] SortableObject, order [] int, range_ Range, comp CompFunction, first int, second int) {
	compareResult := comp(input[range_.start+first], input[range_.start+second])
	if compareResult > 0 || (order[first] > order[second] && compareResult == 0) {
		swap := input[range_.start+first]
		swapOrder := order[first]

		input[range_.start+first] = input[range_.start+second]
		input[range_.start+second] = swap

		order[first] = order[second]
		order[second] = swapOrder
	}
}

// rotate the values in an array ([0 1 2 3] becomes [1 2 3 0] if we rotate by 1)
// this assumes that 0 <= amount <= range.length()
func rotateInput(input [] SortableObject, amount int, range_ Range, useCache bool) {
	if range_.Length() == 0 {
		return
	}
	// TODO: HERE
}

// merge two ranges from one array and save the results into a different array
func mergeInto(input [] SortableObject, firstRange Range, secondRange Range, comp CompFunction, into [] SortableObject, atIndex int) {

}

// merge operation using an external buffer
func externalMarge(input [] SortableObject, firstRange Range, secondRange Range, comp CompFunction) {

}

// merge operation using an internal buffer
func internalMerge(input [] SortableObject, firstRange Range, secondRange Range, comp CompFunction, buffer Range) {

}

// merge operation without a buffer
func inPlaceMarge(input [] SortableObject, firstRange Range, secondRange Range, comp CompFunction) {

}

// ***************** MERGE SORT ******************************
var input = []int{124, 1, 55, 12, 12, 0, 5, 99}
var cache = make([]int, len(input)/2+1)

func mergeSort(input []int) {
	if len(input) < 2 {
		return
	}

	mid := len(input) / 2

	mergeSort(input[:mid])
	mergeSort(input[mid:])

	if input[mid-1] <= input[mid] {
		return
	}
	copy(cache, input[:mid])
	l, r := 0, mid
	for i := 0; ; i++ {
		if cache[l] <= input[r] {
			input[i] = cache[i]
			l++
			if l == mid {
				break
			}
		} else {
			input[i] = input[r]
			r++
			if r == len(input) {
				copy(input[i+1:], cache[l:mid])
				break
			}
		}
	}
	return
}

func main() {

}
