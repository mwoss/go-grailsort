package wikisort

// Structure to represent ranges within the array
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

func (i *Iterator) nextRange() Range {
	start := i.decimal
	i.decimal += i.decimalStep
	i.numerator += i.numeratorStep

	if i.numerator >= i.denominator {
		i.numerator -= i.denominator
		i.decimal++
	}
	return Range{start: start, end: i.decimal}
}

func (i Iterator) isFinished() bool {
	return i.decimal >= i.size
}

func (i *Iterator) isNextLevel() bool {
	i.decimalStep += i.decimalStep
	i.numeratorStep += i.numeratorStep

	if i.numeratorStep >= i.denominator {
		i.numeratorStep -= i.denominator
		i.decimalStep++
	}
	return i.decimalStep < i.size
}


