package password

import "crypto/sha512"

func Hash(s []byte) [sha512.Size]byte {
	return sha512.Sum512(s)
}

func check(password string, req string) bool {
	return Hash([]byte(password)) == Hash([]byte(req))
}
