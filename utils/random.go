package utils

import "math/rand"

/*** Code for lib ***/
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 52 letter, 52 = 110100b => keep 6 smallest bit, zer0 the rest
	letterIdxMask = 1<<letterIdxBits - 1 // 1b move right 6 pos -> 1000000, -1 and turn into 0111111
)
const (
	letterIdxMax = 63 / letterIdxBits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 { // remain = 0 -> init new random number and assign to cache and restart remain
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		// shift cache 6 bit to the right
		cache = cache >> letterIdxBits
		remain--
	}
	return string(b)
}

/*** End code for lib ***/
