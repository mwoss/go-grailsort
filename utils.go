package wikisort

// 63 -> 32, 64 -> 64, etc.
// this comes from Hacker's Delight
func FloorPowerOfTwo(value int) int {
	x := value
	x = x | (x >> 1)
	x = x | (x >> 2)
	x = x | (x >> 4)
	x = x | (x >> 8)
	x = x | (x >> 16)
	return x - (x >> 1)
}
