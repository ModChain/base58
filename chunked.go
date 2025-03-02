package base58

import (
	"fmt"
)

var (
	encodedBlockSizes = []int{0, 2, 3, 5, 6, 7, 9, 10, 11}
	decodedBlockSizes = []int{0, -1, 1, 2, -1, 3, 4, 5, -1, 6, 7, 8}
	fullBlockSize     = 8 // (length of encodedBlockSizes) - 1
	fullEncodedSize   = 11
)

// encodeBlock encodes a single chunked block
func (e *Encoding) encodeBlock(block []byte) string {
	size := len(block)
	if size < 1 || size > fullBlockSize {
		panic("invalid block data length")
	}

	// Convert up to 8 bytes into a 64-bit big-endian integer.
	num := uint64(0)
	for _, b := range block {
		num = (num << 8) | uint64(b)
	}

	// How many base58 chars we should get from this chunk:
	encodedLen := encodedBlockSizes[size]
	res := make([]byte, encodedLen)

	// Convert num into base58 from the right to the left.
	i := encodedLen - 1
	for num > 0 {
		rem := num % 58
		num /= 58
		res[i] = e.encode[rem]
		i--
	}

	// Fill the remainder with alphabet 0
	for j := i; j >= 0; j-- {
		res[j] = e.encode[0]
	}

	return string(res)
}

func (e *Encoding) EncodeChunked(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	var result []byte

	fullBlockCount := len(data) / fullBlockSize
	lastBlockSize := len(data) % fullBlockSize

	// Encode each full block
	for i := 0; i < fullBlockCount; i++ {
		block := data[i*fullBlockSize : (i+1)*fullBlockSize]
		result = append(result, e.encodeBlock(block)...)
	}

	// Encode the remainder
	if lastBlockSize > 0 {
		block := data[fullBlockCount*fullBlockSize:]
		result = append(result, e.encodeBlock(block)...)
	}

	return string(result)
}

func (e *Encoding) decodeBlock(block string) ([]byte, error) {
	size := len(block)
	if size < 1 || size > fullEncodedSize {
		return nil, ErrInvalidBlockLength
	}

	// Determine how many raw bytes we should get
	rawSize := decodedBlockSizes[size]

	// Convert base58 string -> number
	var resNum uint64
	for _, r := range block {
		if r > 127 {
			return nil, ErrNonAscii
		}
		if e.decode[r] == -1 {
			return nil, fmt.Errorf("%w (%q)", ErrBadDigit, r)
		}
		idx := e.decode[r]
		resNum = (resNum * 58) + uint64(idx)
	}

	// Output array
	result := make([]byte, rawSize)
	for n := range result {
		result[rawSize-1-n] = byte((resNum >> (n * 8)) & 0xff)
	}

	return result, nil
}

func (e *Encoding) DecodeChunked(encoded string) ([]byte, error) {
	if len(encoded) == 0 {
		return []byte{}, nil
	}

	fullBlockCount := len(encoded) / fullEncodedSize
	lastBlockSize := len(encoded) % fullEncodedSize

	var result []byte

	// Decode full blocks
	for i := 0; i < fullBlockCount; i++ {
		block := encoded[i*fullEncodedSize : (i+1)*fullEncodedSize]
		decoded, err := e.decodeBlock(block)
		if err != nil {
			return nil, err
		}
		result = append(result, decoded...)
	}

	// Decode the remainder
	if lastBlockSize > 0 {
		block := encoded[fullBlockCount*fullEncodedSize:]
		decoded, err := e.decodeBlock(block)
		if err != nil {
			return nil, err
		}
		result = append(result, decoded...)
	}

	return result, nil
}
