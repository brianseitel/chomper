package chomper

// size of a bit set
const wordSize = uint(64)

// wordSizeLowerBound is the log of wordSize
const wordSizeLowerBound = uint(6)

// Bitset thingy
type Bitset struct {
	length uint
	bits   []uint64
}

// New bitset generator
func New(i int) *Bitset {
	length := uint(i)
	bs := &Bitset{
		length: length,
		bits:   make([]uint64, wordsNeeded(length)),
	}

	return bs
}

// Length of the bitset
func (bs *Bitset) Length() uint {
	return bs.length
}

// Get the bit for this pos
func (bs *Bitset) Get(i int) bool {
	pos := uint(i)
	if pos >= bs.length {
		return false
	}
	return bs.bits[pos>>wordSizeLowerBound]&(1<<(pos&(wordSize-1))) != 0
}

// Set bit i to 1
func (bs *Bitset) Set(pos uint) *Bitset {
	bs.extendSetMaybe(pos)
	bs.bits[pos>>wordSizeLowerBound] |= 1 << (pos & (wordSize - 1))
	return bs
}

// Clear bit i to 0
func (bs *Bitset) Clear(i int) *Bitset {
	pos := uint(i)
	if pos >= bs.length {
		return bs
	}
	bs.bits[pos>>wordSizeLowerBound] &^= 1 << (pos & (wordSize - 1))
	return bs
}

// Count (number of set bits)
func (bs *Bitset) Count() uint {
	if bs != nil && bs.bits != nil {
		return uint(popcntSliceGo(bs.bits))
	}
	return 0
}

// extendSetMaybe adds additional words to incorporate new bits if needed
func (bs *Bitset) extendSetMaybe(i uint) {
	if i >= bs.length { // if we need more bits, make 'em
		nsize := wordsNeeded(i + 1)
		if bs.bits == nil {
			bs.bits = make([]uint64, nsize)
		} else if cap(bs.bits) >= nsize {
			bs.bits = bs.bits[:nsize] // fast resize
		} else if len(bs.bits) < nsize {
			newset := make([]uint64, nsize, 2*nsize) // increase capacity 2x
			copy(newset, bs.bits)
			bs.bits = newset
		}
		bs.length = i + 1
	}
}

// wordsNeeded calculates the number of words needed for i bits
func wordsNeeded(i uint) int {
	cap := ^uint(0) // total possible number of bits
	if i > (cap - wordSize + 1) {
		return int(cap >> wordSizeLowerBound)
	}
	return int((i + (wordSize - 1)) >> wordSizeLowerBound)
}

// bit population count, take from
// https://code.google.com/p/go/issues/detail?id=4988#c11
// credit: https://code.google.com/u/arnehormann/
func popcount(x uint64) (n uint64) {
	x -= (x >> 1) & 0x5555555555555555
	x = (x>>2)&0x3333333333333333 + x&0x3333333333333333
	x += x >> 4
	x &= 0x0f0f0f0f0f0f0f0f
	x *= 0x0101010101010101
	return x >> 56
}

func popcntSliceGo(s []uint64) uint64 {
	cnt := uint64(0)
	for _, x := range s {
		cnt += popcount(x)
	}
	return cnt
}
