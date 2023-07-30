package utils

import "encoding/base64"

var base64Encoding = base64.StdEncoding

func DecodeBase64(data []byte) ([]byte, error) {
	res := make([]byte, base64Encoding.DecodedLen(len(data)))
	_, err := base64Encoding.Decode(res, data)
	if err != nil {
		return nil, err
	}
	return res, nil
}
