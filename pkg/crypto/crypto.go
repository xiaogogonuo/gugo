package crypto

import (
	"crypto/sha1"
	"unsafe"
)

func SHA1Encrypt2Byte(b *[]byte) []byte {
	h := sha1.New()
	h.Write(*b)
	return h.Sum(nil)
}

func SHA1Encrypt2String(b *[]byte) string {
	r := SHA1Encrypt2Byte(b)
	return *(*string)(unsafe.Pointer(&r))
}
