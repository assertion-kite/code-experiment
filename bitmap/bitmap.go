package bitmap

// Bitmap struct
type Bitmap struct {
	bits []uint64
}

// NewBitmap creates a new Bitmap with a given size
func NewBitmap(size int) *Bitmap {
	return &Bitmap{
		bits: make([]uint64, (size+63)/64),
	}
}

func (b *Bitmap) Set(index int) {
	if index < 0 {
		return
	}
	b.bits[index/64] |= 1 << (index % 64)
}

func (b *Bitmap) Test(index int) bool {
	if index < 0 {
		return false
	}
	return b.bits[index/64]&(1<<(index%64)) != 0
}
