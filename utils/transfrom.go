package utils

import "encoding/hex"

func Byte32toStr(byte32 [32]byte) string {
	if byte32 == [32]byte{} {
		return ""
	}
	
	return hex.EncodeToString(byte32[:])
}