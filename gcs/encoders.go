package gcs

import "github.com/golang/snappy"

type transformer func(in []byte) (out []byte, err error)

//noinspection GoUnusedGlobalVariable
var (
	TransformerNone = func(in []byte) (out []byte, err error) {
		return in, nil
	}

	TransformerSnappyEncode = func(in []byte) (out []byte, err error) {
		return snappy.Encode(nil, in), nil
	}
	TransformerSnappyDecode = func(in []byte) (out []byte, err error) {
		return snappy.Decode(nil, in)
	}
)
