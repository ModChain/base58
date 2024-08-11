package base58

// Encode will encode the provided byte array into a base58 encoded string using
// the current encoding.
func (e *Encoding) Encode(bin []byte) string {
	size := len(bin)

	if size == 0 {
		return ""
	}

	zcount := 0
	for zcount < size && bin[zcount] == 0 {
		zcount += 1
	}

	// It is crucial to make this as short as possible, especially for
	// the usual case of bitcoin addrs
	size = zcount + (size-zcount)*555/406 + 1 // This is an integer simplification of ceil(log(256)/log(58))

	out := make([]byte, size)

	var i, high int
	var carry uint32

	high = size - 1
	for _, b := range bin {
		i = size - 1
		for carry = uint32(b); i > high || carry != 0; i-- {
			carry = carry + 256*uint32(out[i])
			out[i] = byte(carry % 58)
			carry /= 58
		}
		high = i
	}

	// Determine the additional "zero-gap" in the buffer (aside from zcount)
	for i = zcount; i < size && out[i] == 0; i++ {
	}

	// Now encode the values with actual digits in-place
	val := out[i-zcount:]
	size = len(val)
	ent := e.encode // encode table
	for i = 0; i < size; i++ {
		out[i] = ent[val[i]]
	}

	return string(out[:size])
}

// EncodedLen returns the encoded len for a given buffer, checking it for initial zeroes. It can be higher than needed.
func EncodedLen(bin []byte) (int, int) {
	size := len(bin)
	if size == 0 {
		return 0, 0
	}
	zcount := 0
	for zcount < size && bin[zcount] == 0 {
		zcount += 1
	}

	return zcount, zcount + (size-zcount)*555/406 + 1 // This is an integer simplification of ceil(log(256)/log(58))
}

// EncodedMaxLen return the maximum encoded len for an input of ln bytes. If it starts with zeroes the actual len will be shorter
func EncodedMaxLen(ln int) int {
	if ln == 0 {
		return 0
	}
	return ln*555/406 + 1
}

// EncodeTo will encoded the provided byte array to the given buffer if it is large enough, or allocate
// a new one if not.
func (e *Encoding) EncodeTo(dst, src []byte) []byte {
	zcount, size := EncodedLen(src)
	if size == 0 {
		if dst == nil {
			return nil
		}
		return dst[:0]
	}

	if cap(dst) >= size {
		dst = dst[:size]
	} else {
		dst = make([]byte, size)
	}

	var i, high int
	var carry uint32

	high = size - 1
	for _, b := range src {
		i = size - 1
		for carry = uint32(b); i > high || carry != 0; i-- {
			carry += uint32(dst[i]) << 8
			dst[i] = byte(carry % 58)
			carry /= 58
		}
		high = i
	}

	// Determine the additional "zero-gap" in the buffer (aside from zcount)
	for i = zcount; i < size && dst[i] == 0; i++ {
	}

	// Now encode the values with actual digits in-place
	dst = dst[i-zcount:]
	ent := e.encode // encode table
	for i = range dst {
		dst[i] = ent[dst[i]]
	}

	return dst
}
