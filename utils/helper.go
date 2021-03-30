package utils

import (
	"encoding/hex"
)

func DecodeLogin(userid string)[]byte{
	//Decodes userid to []byte ...

	dst := make([]byte, hex.DecodedLen(len(userid)))
	hex.Decode(dst ,[]byte(userid))
	return dst
}

func EncodeLogin(login, pass string)[]byte{
	//Encodes userid to hex ...

	dst := make([]byte, hex.EncodedLen(len(login + "." + pass)))
	hex.Encode(dst, []byte(login + "." + pass))
	return dst
}