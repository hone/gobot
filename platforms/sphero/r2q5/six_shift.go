package r2q5

func highestOrderBitIndex(num uint16) uint16 {
	ret := uint16(0)
	n := num
	for n > 0 {
		n = n >> 1
		ret = ret + 1
	}

	return ret
}

func min(a, b uint16) uint16{
	if a < b {
		return a
	}
	return b
}

func DecodeSixShifted(bytes []byte) (num int16, err error) {
	left := bytes[0]
	right := bytes[1]

	if left == 0 && right == 0 {
		return 0, nil
	}

	sign := 0x80 & left >> 7

	if left & 0x0f == 0x0f {
		if sign == 1 {
			return -1, nil
		} else {
			return 1, nil
		}
	}

	A := 2 & left >> 1
	B := 1 & left
	C := 0x80 & right >> 7

	shift := A << 2 + B << 1 + C
	val := int16((0x80 | right) >> (6 - shift))
	if (sign == 1) {
		val = val * -1
	}
	return val, nil
}

func EncodeSixShifted(num int16) (bytes []byte, err error) {
	if num == -1 {
		return []byte{0xbf,0x80}, nil
	}
	if num == 0 {
		return []byte{0x00,0x00}, nil
	}
	if num == 1 {
		return []byte{0x3f,0x80}, nil
	}

	neg := num < 0

	var unum uint16
	if neg {
		unum = uint16(-num)
	} else {
		unum = uint16(num)
	}

	hob := highestOrderBitIndex(unum)
	unshift := min(8 - hob, 6)
	shift := 6 - unshift

	unum = unum << unshift

	if neg {
		unum = 0x8000 | unum
	}
	unum = 0x4000 | unum

	A := (0x04 & shift) >> 2
	B := (0x02 & shift) >> 1
	C := 0x01 & shift

	if A == 1 {
		unum |= 1 << 9
	} else {
		unum &= ^(uint16(1) << 9)
	}

	if B == 1 {
		unum |= 1 << 8
	} else {
		unum &= ^(uint16(1) << 8)
	}

	if C == 1 {
		unum |= 1 << 7
	} else {
		unum &= ^(uint16(1) << 7)
	}

	lsb := byte(0x00FF & unum)
	msb := byte((0xFF00 & unum) >> 8)

	return []byte{msb,lsb}, nil
}
