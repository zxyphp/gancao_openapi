package tool

import (
	"math/rand"
	"time"
)

const (
	letterBytes   = "ABCDEFGHIJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandStr(iLength int) string {
	if iLength == 0 {
		iLength = 4
	}

	// Shuffle the letters
	shuffledLetters := shuffleString(letterBytes)

	// Prepare output buffer
	var outBuf []rune
	// Generate the first character (it should not be '0')
	for {
		r := rune(shuffledLetters[rand.Intn(len(shuffledLetters))])
		if r != '0' {
			outBuf = append(outBuf, r)
			break
		}
	}

	// Generate the rest of the characters
	for len(outBuf) < iLength {
		r := rune(shuffledLetters[rand.Intn(len(shuffledLetters))])
		if len(outBuf) == 0 || r != outBuf[len(outBuf)-1] {
			outBuf = append(outBuf, r)
		}
	}

	return string(outBuf)
}

// shuffleString shuffles the characters in a string
func shuffleString(s string) string {
	r := []rune(s)
	for i := range r {
		j := rand.Intn(i + 1)
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
