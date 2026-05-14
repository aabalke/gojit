package arm64

import (
	"math"
	"math/bits"
)

// https://kddnewton.com/2022/08/11/aarch64-bitmask-immediates.html

type Immediate struct {
	n, imms, immr uint32
}

func BuildImmediate(imm uint64) Immediate {

	if imm == 0 || imm == math.MaxUint64 {
		panic("imm with v == 0 || v == MAX 64")
	}

	size := uint64(64)

	for {
		size >>= 1
		mask := uint64(1 << size) - 1

		if imm & mask != (imm >> size) & mask {
			size <<= 1
			break
		}

		if size <= 2 {
			break
		}
	}

	var trailingOnes, leftRots uint32

	mask := uint64(math.MaxUint64) >> (64 - size)
	imm &= mask

	if _isShiftedMask(imm) {
		leftRots = uint32(bits.TrailingZeros64(imm))
		trailingOnes = uint32(bits.TrailingZeros64(^(imm >> uint64(leftRots))))
	} else {
		imm |= ^mask
		if !_isShiftedMask(^imm) {
			panic("invalid imm v")
		}

		leadingOnes := uint32(bits.LeadingZeros64(^imm))

		leftRots = 64 - leadingOnes
		trailingOnes = leadingOnes + uint32(bits.TrailingZeros64(^imm)) - uint32(64 - size)
	}

	immr := (size - uint64(leftRots)) & (size - 1)
	imms := (^(size - 1) << 1) | uint64(trailingOnes - 1)
	n := ((imms >> 6) & 1) ^ 1

	return Immediate {
		n: uint32(n),
		imms: uint32(imms & 0x3f),
		immr: uint32(immr & 0x3f),
	}
}

func _isMask(imm uint64) bool {
	return (imm + 1) & imm == 0
}

func _isShiftedMask(imm uint64) bool {
	return _isMask((imm - 1) | imm)
}
