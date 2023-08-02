package utils

import "encoding/base64"

var base64Encoding = base64.StdEncoding.WithPadding(base64.NoPadding)

func DecodeBase64(data []byte) ([]byte, error) {
	res := make([]byte, base64Encoding.DecodedLen(len(data)))
	_, err := base64Encoding.Decode(res, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func EncodeBase64(data []byte) []byte {
	res := make([]byte, base64Encoding.EncodedLen(len(data)))
	base64Encoding.Encode(res, data)
	return res
}
