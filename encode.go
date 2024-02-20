package base58

func (e *Encoding) Encode(bin []byte) string {
	size := len(bin)

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
	for i = 0; i < size; i++ {
		out[i] = e.encode[val[i]]
	}

	return string(out[:size])
}
