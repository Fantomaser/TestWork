package main

import (
	"fmt"
	"math/rand"

)

func Encrypt(str []byte) ([]byte, []byte) {

	btMask := criptKey()

	for it := 0; len(UserDB) != 0; it++ {
		if UserDB[it].Key == string(btMask) {
			btMask = criptKey()
			it = 0
		}

		if it == (len(UserDB) - 1) {
			break
		}

	}

	for i, _ := range []byte(str) {
		str[i] = str[i] ^ btMask[i%len(btMask)]
	}

	return str, btMask
}

func criptKey() []byte {

	btMask := make([]byte, 8)

	for i, _ := range btMask {
		btMask[i] = byte(rand.Int63() % 126)
	}

	return btMask
}

func Decrypt(key []byte) ([]byte, error) {

	var Userid int = -1

	for i, _ := range UserDB {
		if UserDB[i].Key == string(key) {
			Userid = i
			break
		}
	}

	if Userid == -1 {
		return nil, fmt.Errorf("deny")
	}

	str := []byte(UserDB[Userid].Message)

	for i, _ := range str {
		str[i] = str[i] ^ key[i%len(key)]
	}

	return str, nil
}
