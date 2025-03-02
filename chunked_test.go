package base58_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/ModChain/base58"
)

func TestChunked(t *testing.T) {
	vectors := []string{
		"00",
		"00000000",
		"01234567abcbcdef",
		"01234567abcbcdef01234567abcbcdef01234567abcbcdef01234567abcbcdef",
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
