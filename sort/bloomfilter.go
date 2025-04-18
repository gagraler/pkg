package sort

import (
	"hash/fnv"
	"math"
)

// 布隆过滤器

type BloomFilter struct {
	bitset *BitMap
	k      int // 哈希函数个数
	size   int
}

func NewBloomFilter(n int, fp float64) *BloomFilter {
	m := int(math.Ceil(float64(n) * math.Log(fp) / math.Log(1/math.Pow(2, math.Log(2)))))
	k := int(math.Round(float64(m) / float64(n) * math.Log(2)))

	return &BloomFilter{
		bitset: NewBitMap(m),
		k:      k,
		size:   m,
	}
}

func (bf *BloomFilter) Add(data []byte) {
	h := bf.hashes(data)
	for _, v := range h {
		bf.bitset.Set(v % bf.size)
	}
}

func (bf *BloomFilter) Contains(data []byte) bool {
	h := bf.hashes(data)
	for _, v := range h {
		if !bf.bitset.Get(v % bf.size) {
			return false
		}
	}
	return true
}

// 多哈希函数模拟（双重哈希法）
func (bf *BloomFilter) hashes(data []byte) []int {
	h1 := fnv.New64a()
	h1.Write(data)
	a := int(h1.Sum64())

	h2 := fnv.New64()
	h2.Write(data)
	b := int(h2.Sum64())

	hashes := make([]int, bf.k)
	for i := 0; i < bf.k; i++ {
		hashes[i] = a + i*b
	}
	return hashes
}
