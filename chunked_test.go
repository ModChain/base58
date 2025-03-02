package base58_test

import (
	"bytes"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ModChain/base58"
)

func TestChunked(t *testing.T) {
	vectors := []string{
		"00",
		"00000000",
		"01234567abcbcdef",
		"01234567abcbcdef01234567abcbcdef01234567abcbcdef01234567abcbcdef",
		"0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f",
	}

	// extra test vectors
	for n := 0; n < 32; n++ {
		vectors = append(vectors,
			strings.Repeat("00", n),
			strings.Repeat("ff", n),
			strings.Repeat("55", n),
		)
	}

	for _, vechex := range vectors {
		vec := must(hex.DecodeString(vechex))
		enc := base58.Bitcoin.EncodeChunked(vec)
		dec, err := base58.Bitcoin.DecodeChunked(enc)
		if err != nil {
			t.Errorf("failed to decode encoded string: %s", err)
			continue
		}
		if !bytes.Equal(vec, dec) {
			t.Errorf("failed to decode back to input: %x → %s → %x", vec, enc, dec)
		}
	}
}

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
