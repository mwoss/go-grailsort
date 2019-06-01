package wikisort

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

func (w WikiSorter) Sort(input []SortableObject, comperatorFun func(elem1 SortableObject, elem2 SortableObject) bool) {

}
var input = []int{124, 1, 55, 12, 12,0, 5, 99}
var cache = make([]int, len(input) / 2 + 1)

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
		if cache[l] <= input[r]{
			input[i] = cache[i]
			l++
			if l == mid {
				break
			}
		} else {
			input[i] = input[r]
			r++
			if r == len(input){
				copy(input[i+1:], cache[l:mid])
				break
			}
		}
	}
	return
}

func main() {

}