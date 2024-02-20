package base58

import "errors"

var (
	ErrZeroLength = errors.New("base58: cannot decode zero length string")
	ErrNonAscii   = errors.New("base58: cannot decode non-ASCII input")
	ErrBadDigit   = errors.New("base58: cannot decode unsupported digit")
)
