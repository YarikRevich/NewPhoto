package utils

import (
	"github.com/YarikRevich/NewPhoto/log"
	"encoding/hex"
)

func DecodeLogin(userid string) []byte {
	//Decodes userid to []byte ...

	dst := make([]byte, hex.DecodedLen(len(userid)))
	if _, err := hex.Decode(dst, []byte(userid)); err != nil {
		log.Logger.UsingErrorLogFile().CFatalln("DecodeLogin", err)
	}
	return dst
}

func EncodeLogin(login, pass string) []byte {
	//Encodes userid to hex ...

	dst := make([]byte, hex.EncodedLen(len(login+"."+pass)))
	hex.Encode(dst, []byte(login+"."+pass))
	return dst
}
