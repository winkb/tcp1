package types

type NumInt interface {
	int | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

type NumFloat interface {
	float32 | float64
}

type Number interface {
	NumInt | NumFloat
}
