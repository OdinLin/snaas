package generate

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"time"

	"golang.org/x/crypto/scrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ~!$%^&*()_+{}:\"|<>?`-=[];'\\,./"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits

	saltLen      = 256
	scryptN      = 32768
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 256
)

// EncryptPassword encrypts the given password with scrypt and the pre-defined
// security parameters.
func EncryptPassword(pw, salt []byte) ([]byte, error) {
	return scrypt.Key(pw, salt, scryptN, scryptR, scryptP, scryptKeyLen)
}

// UUID generates a random UUID according to RFC 4122.
func UUID() (string, error) {
	uuid := make([]byte, 16)

	n, err := io.ReadFull(crand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}

	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40

	return fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:],
	), nil
}

// Salt returns a secure byte string to use for password securing.
func Salt() ([]byte, error) {
	var (
		salt = make([]byte, saltLen)
	)

	_, err := crand.Read(salt)
	return salt, err
}

// RandomBytes returns a generated bytes with the provided length.
//
// Solution based on SO thread: http://stackoverflow.com/a/31832326/1590256
func RandomBytes(src rand.Source, n int) []byte {
	b := make([]byte, n)

	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return b
}

// RandomString returns a genrated string with the provided length.
func RandomString(n int) string {
	return string(RandomBytes(rand.NewSource(time.Now().UnixNano()), n))
}
