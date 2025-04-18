package sort

// BitMap实现

type BitMap struct {
	data []uint64
}

func NewBitMap(size int) *BitMap {
	return &BitMap{
		data: make([]uint64, (size+63)/64), // 每个 uint64 表示 64 个 bit
	}
}

func (b *BitMap) Set(pos int) {
	idx := pos / 64
	off := pos % 64
	b.data[idx] |= 1 << off
}

func (b *BitMap) Clear(pos int) {
	idx := pos / 64
	off := pos % 64
	b.data[idx] &^= 1 << off
}

func (b *BitMap) Get(pos int) bool {
	idx := pos / 64
	off := pos % 64
	return b.data[idx]&(1<<off) != 0
}

func (b *BitMap) List() []int {
	var result []int
	for i, v := range b.data {
		for j := 0; j < 64; j++ {
			if v&(1<<j) != 0 {
				result = append(result, i*64+j)
			}
		}
	}
	return result
}

func (b *BitMap) Size() int {
	return len(b.data) * 64
}

func (b *BitMap) Count() int {
	count := 0
	for _, v := range b.data {
		count += popcount(v)
	}
	return count
}

func popcount(x uint64) int {
	count := 0
	for x > 0 {
		count++
		x &= x - 1 // 清除最低位的 1
	}
	return count
}
