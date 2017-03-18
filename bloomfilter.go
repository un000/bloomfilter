package bloomfilter

import (
	"github.com/spaolacci/murmur3"
	"github.com/willf/bitset"
	"hash"
	"hash/fnv"
	"math"
)

type Filter struct {
	m uint // storage size
	k uint // hashes count

	bits *bitset.BitSet

	FirstHasher  hash.Hash64 // the first hasher
	SecondHasher hash.Hash64 // the second hasher

	hashes []uint // computed hash storage
}

// n - size
func New(n uint) *Filter {
	f := &Filter{
		FirstHasher:  fnv.New64(),
		SecondHasher: murmur3.New64(),
	}

	p := float64(0.5)
	e := float64(0.001)

	f.m = f.predictM(n, p, e)
	f.k = f.predictK(e)

	f.bits = bitset.New(uint(f.m))
	f.hashes = make([]uint, f.k)

	return f
}

func (f *Filter) Insert(value []byte) {
	f.computeHashes(value)

	for i := range f.hashes {
		f.bits.Set(f.hashes[i])
	}
}

func (f *Filter) Has(value []byte) bool {
	f.computeHashes(value)

	for i := range f.hashes {
		if !f.bits.Test(f.hashes[i]) {
			return false
		}
	}

	return true
}

// e - error
func (f *Filter) predictK(e float64) uint {
	return uint(math.Ceil(math.Log2(1 / e)))
}

// n - count
// p - fill ratio
// e - error
func (f *Filter) predictM(n uint, p float64, e float64) uint {
	return uint(math.Ceil(float64(n) / ((math.Log(p) * math.Log(1-p)) / math.Abs(math.Log(e)))))
}

// compute hashes
func (f *Filter) computeHashes(value []byte) {
	f.FirstHasher.Reset()
	f.SecondHasher.Reset()
	f.FirstHasher.Write(value)
	f.SecondHasher.Write(value)

	h1 := f.FirstHasher.Sum64()
	h2 := f.SecondHasher.Sum64()

	for i := uint(0); i < f.k; i++ {
		g := (uint(h1) + uint(h2)*(i+1)) % f.m
		f.hashes[i] = g
	}
}
