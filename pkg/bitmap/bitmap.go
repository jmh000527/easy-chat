package bitmap

// Bitmap 是一个位图结构，用于高效地表示和操作大量的布尔值。
type Bitmap struct {
	bits []byte // 存储位图的字节数组
	size int    // 位图的总位数
}

// NewBitmap 创建并返回一个新的 Bitmap 实例。
// 参数 size 指定了位图的大小，以位数计。
// 如果 size 小于等于 0，将默认为 250。
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

// Set 将给定 id 对应的位设置为 1。
// 参数 id 是一个唯一标识，用于定位位图中的位。
func (bm *Bitmap) Set(id string) {
	// 计算 id 在位图中的索引位置
	// 计算ID在哪个bit
	idx := hash(id) % bm.size
	// 根据bit计算哪个字节
	byteIndex := idx / 8
	// 计算在该字节中该bit的偏移量
	bitIndex := idx % 8
	// 通过位运算将对应的位设置为 1
	bm.bits[byteIndex] |= 1 << bitIndex
}

// IsSet 检查给定 id 对应的位是否被设置为 1。
// 参数 id 是一个唯一标识，用于定位位图中的位。
// 返回值表示位的状态：true 表示位为 1，false 表示位为 0。
func (bm *Bitmap) IsSet(id string) bool {
	// 计算 id 在位图中的索引位置
	idx := hash(id) % bm.size
	byteIndex := idx / 8
	bitIndex := idx % 8
	// 通过位运算检查对应的位是否为 1
	return bm.bits[byteIndex]&(1<<bitIndex) != 0
}

// Export 返回位图的字节数组表示。
func (bm *Bitmap) Export() []byte {
	return bm.bits
}

// Load 从字节数组创建并返回一个 Bitmap 实例。
// 参数 bits 是位图的字节数组表示。
// 如果 bits 为空，将返回一个大小为 0 的位图。
func Load(bits []byte) *Bitmap {
	if len(bits) == 0 {
		return NewBitmap(0)
	}
	return &Bitmap{
		bits: bits,
		size: len(bits) * 8,
	}
}

// hash 计算字符串 id 的哈希值。
// 该哈希函数用于将 id 映射到位图的索引位置。
// 参数 id 是需要计算哈希值的字符串。
// 返回值是 id 的哈希值，映射到位图索引的范围。
func hash(id string) int {
	// 使用 BKDR 哈希算法
	seed := 131313 // 31 131 1313 13131 131313, etc
	hash := 0
	for _, c := range id {
		hash = hash*seed + int(c)
	}
	// 保证哈希值在 int 的正范围内
	return hash & 0x7FFFFFFF
}
