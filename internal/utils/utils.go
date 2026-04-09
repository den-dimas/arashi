package utils

import "math"

func DecodeFloat(reg1, reg2 uint16) float32 {
	bits := uint32(reg1)<<16 | uint32(reg2)

	return math.Float32frombits(bits)
}

func EncodeFloat(value float32) []uint16 {
	bits := math.Float32bits(value)

	return []uint16{uint16(bits >> 16), uint16(bits & 0xffff)}
}
