package gojit

import (
	"math/bits"
)

// https://dougallj.wordpress.com/2021/10/30/bit-twiddling-optimising-aarch64-logical-immediate-encoding-and-decoding/
// CC0 Dougall Johnson
// original did not support 32bit only patterns

type Immediate struct {
	n, imms, immr uint32
	sf            bool
}

func EncodeImm(val uint64, sf bool) Immediate {
	if sf {
		imm := val
		if imm == 0 || ^imm == 0 {
			panic("imm with v == 0 || v == MAX 64")
		}

		var (
			rotation   = bits.TrailingZeros64(imm & (imm + 1))
			normalized = bits.RotateLeft64(imm, -(rotation & 63))
			zeros      = bits.LeadingZeros64(normalized)
			ones       = bits.TrailingZeros64(^normalized)
			size       = zeros + ones
		)

		if bits.RotateLeft64(imm, -(size&63)) != imm {
			panic("invalid imm: non-repeating pattern")
		}

		return Immediate{
			n:    uint32(size >> 6),
			immr: uint32((-rotation) & (size - 1)),
			imms: uint32((-(size << 1) | (ones - 1)) & 0x3f),
		}
	}

	imm := uint32(val)
	if imm == 0 || ^imm == 0 {
		panic("imm with v == 0 || v == MAX 32")
	}

	var (
		rotation   = bits.TrailingZeros32(imm & (imm + 1))
		normalized = bits.RotateLeft32(imm, -(rotation & 31))
		zeros      = bits.LeadingZeros32(uint32(normalized))
		ones       = bits.TrailingZeros32(uint32(^normalized))
		size       = zeros + ones
	)

	if bits.RotateLeft32(imm, -(size&31)) != imm {
		panic("invalid imm: non-repeating pattern")
	}

	return Immediate{
		n:    uint32(size >> 6),
		immr: uint32((-rotation) & (size - 1)),
		imms: uint32((-(size << 1) | (ones - 1)) & 0x3f),
		sf:   sf,
	}
}

var maskLookup = [6]uint64{
	0xffffffffffffffff, // size = 64
	0x00000000ffffffff, // size = 32
	0x0000ffff0000ffff, // size = 16
	0x00ff00ff00ff00ff, // size = 8
	0x0f0f0f0f0f0f0f0f, // size = 4
	0x3333333333333333, // size = 2
}

func DecodeImmediate(imm Immediate) (uint64, bool) {
	if !imm.sf && imm.n != 0 {
		return 0, false
	}

	pattern := (imm.n << 6) | (^imm.imms & 0x3f)
	// pattern must be a power of two (nonzero)
	if pattern&(pattern-1) == 0 {
		return 0, false
	}

	leadingZeroes := bits.LeadingZeros32(pattern)
	immsMask := uint32(0x7fffffff) >> leadingZeroes
	mask := maskLookup[leadingZeroes-25]
	s := (imm.imms + 1) & immsMask

	v := bits.RotateLeft64(mask^(mask<<s), -int(imm.immr))
	if !imm.sf {
		return v & 0xFFFF_FFFF, true
	}

	return v, true
}
