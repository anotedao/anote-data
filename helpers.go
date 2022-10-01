package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
)

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func EncryptMessage(message string) string {
	byteMsg := []byte(message)
	block, err := aes.NewCipher(conf.Password)
	if err != nil {
		log.Println(err)
	}

	cipherText := make([]byte, aes.BlockSize+len(byteMsg))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Println(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], byteMsg)

	return base64.StdEncoding.EncodeToString(cipherText)
}

func DecryptMessage(message string) string {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		log.Println(err)
	}

	block, err := aes.NewCipher(conf.Password)
	if err != nil {
		log.Println(err)
	}

	if len(cipherText) < aes.BlockSize {
		log.Println(err)
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText)
}
