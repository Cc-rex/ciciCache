package ciciCache

// TODO: encapsulates the slice and ensures that any external access to the data is read-only.

type ByteView struct {
	b []byte
}

func copyByte(b []byte) []byte {
	copyBytes := make([]byte, len(b))
	copy(copyBytes, b)
	return copyBytes
}

func (bv ByteView) Size() int {
	return len(bv.b)
}

// DCopy a deep copy of the slice is returned
func (bv ByteView) DCopy() []byte {
	return copyByte(bv.b)
}

func (bv ByteView) String() string {
	return string(bv.b)
}
