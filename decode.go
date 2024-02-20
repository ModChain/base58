package base58

import (
	"fmt"
)

// Decode will decode the provided string using the current encoding and return
// a byte array of the value, or an error if the input was not valid
func (e *Encoding) Decode(str string) ([]byte, error) {
	if len(str) == 0 {
		return nil, ErrZeroLength
	}

	zero := e.encode[0]
	b58sz := len(str)

	var zcount int
	for i := 0; i < b58sz && str[i] == zero; i++ {
		zcount += 1
	}

	var t, c uint64

	// the 32bit algo stretches the result up to 2 times
	binu := make([]byte, 2*((b58sz*406/555)+1))
	outi := make([]uint32, (b58sz+3)/4)

	for _, r := range str {
		if r > 127 {
			return nil, ErrNonAscii
		}
		if e.decode[r] == -1 {
			return nil, fmt.Errorf("%w (%q)", ErrBadDigit, r)
		}

		c = uint64(e.decode[r])

		for j := len(outi) - 1; j >= 0; j-- {
			t = uint64(outi[j])*58 + c
			c = t >> 32
			outi[j] = uint32(t & 0xffffffff)
		}
	}

	// initial mask depends on b58sz, on further loops it always starts at 24 bits
	mask := (uint(b58sz%4) * 8)
	if mask == 0 {
		mask = 32
	}
	mask -= 8

	outLen := 0
	for j := 0; j < len(outi); j++ {
		for mask < 32 { // loop relies on uint overflow
			binu[outLen] = byte(outi[j] >> mask)
			mask -= 8
			outLen++
		}
		mask = 24
	}

	// find the most significant byte post-decode, if any
	for msb := zcount; msb < len(binu); msb++ {
		if binu[msb] > 0 {
			return binu[msb-zcount : outLen], nil
		}
	}

	// it's all zeroes
	return binu[:outLen], nil
}
