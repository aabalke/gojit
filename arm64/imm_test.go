package arm64

import (
	"fmt"
	"testing"
)

func TestImm(t *testing.T) {

	var tests = []struct {
		Input uint64
		Want Immediate
	}{

	{0x5555555555555555, Immediate{n: 0, immr: 0b000000, imms: 0b111100}},
	{0xaaaaaaaaaaaaaaaa, Immediate{n: 0, immr: 0b000001, imms: 0b111100}},

	{0x1111111111111111, Immediate{n: 0, immr: 0b000000, imms: 0b111000}},
	{0x6666666666666666, Immediate{n: 0, immr: 0b000011, imms: 0b111001}},
	{0xeeeeeeeeeeeeeeee, Immediate{n: 0, immr: 0b000011, imms: 0b111010}},

	{0x0101010101010101, Immediate{n: 0, immr: 0b000000, imms: 0b110000}},
	{0x1818181818181818, Immediate{n: 0, immr: 0b000101, imms: 0b110001}},
	{0xfefefefefefefefe, Immediate{n: 0, immr: 0b000111, imms: 0b110110}},

	{0x0001000100010001, Immediate{n: 0, immr: 0b000000, imms: 0b100000}},
	{0xff8fff8fff8fff8f, Immediate{n: 0, immr: 0b001001, imms: 0b101100}},
	{0xfffefffefffefffe, Immediate{n: 0, immr: 0b001111, imms: 0b101110}},

	{0x0000000100000001, Immediate{n: 0, immr: 0b000000, imms: 0b000000}},
	{0x3fffff003fffff00, Immediate{n: 0, immr: 0b011000, imms: 0b010101}},
	{0xfffffffefffffffe, Immediate{n: 0, immr: 0b011111, imms: 0b011110}},

	{0x0000000000000001, Immediate{n: 1, immr: 0b000000, imms: 0b000000}},
	{0x0000001fffff0000, Immediate{n: 1, immr: 0b110000, imms: 0b010100}},
	{0xfffffffffffffffe, Immediate{n: 1, immr: 0b111111, imms: 0b111110}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%X", tt.Input)
		t.Run(testname, func(t *testing.T) {
			imm := BuildImmediate(tt.Input)
			if !Equal(imm, tt.Want) {
				t.Errorf("got %s, want %s", imm.String(), tt.Want.String())
			}
		})
	}
}

func Equal(a, b Immediate) bool {
	return a.n == b.n && a.immr == b.immr && a.imms == b.imms
}

func (i *Immediate) String() string {
	return fmt.Sprintf("n=%b immr=%b imms=%b", i.n, i.immr, i.imms)
}
