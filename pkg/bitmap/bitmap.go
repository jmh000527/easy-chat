package bitmap

type Bitmap struct {
	bits []byte // 字节数
	size int    // 总共有多少个bit
}

func NewBitmap(size int) *Bitmap {
	if size < 0 {
		panic("size must be greater than zero")
	}
	if size == 0 {
		size = 250
	}
	return &Bitmap{
		bits: make([]byte, size),
		size: size * 8,
	}
}

func (bm *Bitmap) Set(id string) {
	// 计算ID在哪个bit
	idx := hash(id) % bm.size
	// 根据bit计算哪个字节
	byteIndex := idx / 8
	// 计算在该字节中该bit的偏移量
	bitIndex := idx % 8
	bm.bits[byteIndex] |= 1 << bitIndex
}

func (bm *Bitmap) IsSet(id string) bool {
	idx := hash(id) % bm.size
	byteIndex := idx / 8
	bitIndex := idx % 8
	return bm.bits[byteIndex]&(1<<bitIndex) != 0
}

func (bm *Bitmap) Export() []byte {
	return bm.bits
}

func Load(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}
	return &Bitmap{
		bits: bits,
		size: len(bits) * 8,
	}
}

func hash(id string) int {
	// 使用 BKDR 哈希算法
	seed := 131313 // 31 131 1313 13131 131313, etc
	hash := 0
	for _, c := range id {
		hash = hash*seed + int(c)
	}
	return hash & 0x7FFFFFFF
}
