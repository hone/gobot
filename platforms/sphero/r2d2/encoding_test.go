package r2d2

import (
	"testing"
)

func TestEncodeSixShifted(t *testing.T) {
	tests := []struct {
		value   int16
		encoded []byte
	}{
		{-180, []byte{0xc3,0x34} },
		{-160, []byte{0xc3,0x20} },
		{-100, []byte{0xc2,0xc8} },
		{ -64, []byte{0xc2,0x80} },
		{ -32, []byte{0xc2,0x00} },
		{ -16, []byte{0xc1,0x80} },
		{ -15, []byte{0xc1,0x70} },
		{ -10, []byte{0xc1,0x20} },
		{  -4, []byte{0xc0,0x80} },
		{  -3, []byte{0xc0,0x40} },
		{  -2, []byte{0xc0,0x00} },
		{  -1, []byte{0xbf,0x80} },
		{   0, []byte{0x00,0x00} },
		{   1, []byte{0x3f,0x80} },
		{   2, []byte{0x40,0x00} },
		{   3, []byte{0x40,0x40} },
		{  90, []byte{0x42,0xb4} },
		{ 180, []byte{0x43,0x34} },
	}

	for _, tt := range tests {
		actual,_ := EncodeSixShifted(tt.value)
		if actual[0] != tt.encoded[0] || actual[1] != tt.encoded[1] {
			t.Errorf("Expected 0x%x, got 0x%x for data %v.", tt.encoded, actual, tt.value)
		}
	}
}

func TestDecodeSixShifted(t *testing.T) {
	tests := []struct {
		value   int16
		encoded []byte
	}{
		{-180, []byte{0xc3,0x34} },
		{-160, []byte{0xc3,0x20} },
		{-100, []byte{0xc2,0xc8} },
		{ -64, []byte{0xc2,0x80} },
		{ -32, []byte{0xc2,0x00} },
		{ -16, []byte{0xc1,0x80} },
		{ -15, []byte{0xc1,0x70} },
		{ -10, []byte{0xc1,0x20} },
		{  -4, []byte{0xc0,0x80} },
		{  -3, []byte{0xc0,0x40} },
		{  -2, []byte{0xc0,0x00} },
		{  -1, []byte{0xbf,0x80} },
		{   0, []byte{0x00,0x00} },
		{   1, []byte{0x3f,0x80} },
		{   2, []byte{0x40,0x00} },
		{   3, []byte{0x40,0x40} },
		{  90, []byte{0x42,0xb4} },
		{ 180, []byte{0x43,0x34} },
	}

	for _, tt := range tests {
		actual,_ := DecodeSixShifted(tt.encoded)
		if actual != tt.value {
			t.Errorf("Expected %v, got %v for data 0x%x.", tt.value, actual, tt.encoded)
		}
	}
}
