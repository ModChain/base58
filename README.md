[![GoDoc](https://godoc.org/github.com/ModChain/base58?status.svg)](https://godoc.org/github.com/ModChain/base58)

Modified implementation based on `github.com/mr-tron/base58` itself based on https://github.com/trezor/trezor-crypto/blob/master/base58.c

It's nice to have a fast base58 implementation but do we really need to have the slow version in the same lib too?

## Usage

```go
// to decode some base58 string
dec, err := base58.Bitcoin.Decode(in)
if err != nil {
    // handle err
}
// or, to encode:
enc := base58.bitcoin.Encode(dec)
```


