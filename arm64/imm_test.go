package arm64
//
//import (
//	"fmt"
//	"testing"
//)
//
//func TestImm(t *testing.T) {
//	tests := []struct {
//		Input uint64
//		Sf    bool
//		Want  Immediate
//	}{
//		{0x5555555555555555, true, Immediate{n: 0, immr: 0b000000, imms: 0b111100}},
//		{0xaaaaaaaaaaaaaaaa, true, Immediate{n: 0, immr: 0b000001, imms: 0b111100}},
//		{0x1111111111111111, true, Immediate{n: 0, immr: 0b000000, imms: 0b111000}},
//		{0x6666666666666666, true, Immediate{n: 0, immr: 0b000011, imms: 0b111001}},
//		{0xeeeeeeeeeeeeeeee, true, Immediate{n: 0, immr: 0b000011, imms: 0b111010}},
//		{0x0101010101010101, true, Immediate{n: 0, immr: 0b000000, imms: 0b110000}},
//		{0x1818181818181818, true, Immediate{n: 0, immr: 0b000101, imms: 0b110001}},
//		{0xfefefefefefefefe, true, Immediate{n: 0, immr: 0b000111, imms: 0b110110}},
//		{0x0001000100010001, true, Immediate{n: 0, immr: 0b000000, imms: 0b100000}},
//		{0xff8fff8fff8fff8f, true, Immediate{n: 0, immr: 0b001001, imms: 0b101100}},
//		{0xfffefffefffefffe, true, Immediate{n: 0, immr: 0b001111, imms: 0b101110}},
//		{0x0000000100000001, true, Immediate{n: 0, immr: 0b000000, imms: 0b000000}},
//		{0x3fffff003fffff00, true, Immediate{n: 0, immr: 0b011000, imms: 0b010101}},
//		{0xfffffffefffffffe, true, Immediate{n: 0, immr: 0b011111, imms: 0b011110}},
//		{0x0000000000000001, true, Immediate{n: 1, immr: 0b000000, imms: 0b000000}},
//		{0x0000001fffff0000, true, Immediate{n: 1, immr: 0b110000, imms: 0b010100}},
//		{0xfffffffffffffffe, true, Immediate{n: 1, immr: 0b111111, imms: 0b111110}},
//		{0xfffffffc, false, Immediate{n: 0, immr: 30, imms: 29}},
//	}
//
//	for _, tt := range tests {
//		testname := fmt.Sprintf("%X", tt.Input)
//		t.Run(testname, func(t *testing.T) {
//			imm := EncodeImm(tt.Input, tt.Sf)
//			if !Equal(imm, tt.Want) {
//				t.Errorf("got %s, want %s", imm.String(), tt.Want.String())
//			}
//		})
//	}
//}
//
//func Equal(a, b Immediate) bool {
//	return a.n == b.n && a.immr == b.immr && a.imms == b.imms
//}
//
//func (i *Immediate) String() string {
//	return fmt.Sprintf("n=%b immr=%b imms=%b", i.n, i.immr, i.imms)
//}
