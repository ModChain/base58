package base58

var (
	Bitcoin = NewEncoding("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")
	Flickr  = NewEncoding("123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ")
)

type Encoding struct {
	decode [128]int8
	encode [58]byte
}

// NewEncoding returns a encoding structure initialized for decoding/encoding using the passed parameter
//
// It panics if the passed string is not 58 bytes long or isn't valid ASCII.
func NewEncoding(s string) *Encoding {
	if len(s) != 58 {
		panic("base58 encoding must be 58 bytes long")
	}

	ret := &Encoding{}
	copy(ret.encode[:], s)
	for i := range ret.decode {
		ret.decode[i] = -1
	}
	for i, b := range ret.encode {
		if ret.decode[b] != -1 {
			panic("base58 encoding has duplicated symbol")
		}
		ret.decode[b] = int8(i)
	}
	return ret
}
