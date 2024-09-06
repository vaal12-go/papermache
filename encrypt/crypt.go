package encrypt

import (
	"bytes"
	"crypto/rand"
	"fmt"

	"github.com/golang-module/dongle"
)

const AES_CBC_BLOCK = 16

func DecodeCode(code, key string) (string, error) {
	codeBytes := dongle.Decode.FromString(code).ByBase64().ToBytes()
	// fmt.Printf("codeBytes: %v\n", codeBytes)
	plText, err := Decrypt([]byte(codeBytes), []byte(key))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return "", err
	}
	// fmt.Printf("plText: %v\n", plText)
	// fmt.Printf("plText (str): %v\n", string(plText))
	return string(plText), nil
}

func Decrypt(cypherText []byte, key []byte) (plainText []byte, err error) {
	ivBytes := cypherText[:16]
	encrBytes := cypherText[16:]
	stretchedKey := stretchKey(key)
	cipher := dongle.NewCipher()
	cipher.SetMode(dongle.CBC) // CBC、CFB、OFB、CTR、ECB
	//TODO: for some reason dongle PKCS5 padding is not working properly. Have to use custom padding function.
	//TODO: look into source code of dongle to see the problem
	cipher.SetPadding(dongle.No) // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cipher.SetKey(stretchedKey)  // key must be 16, 24 or 32 bytes
	cipher.SetIV(ivBytes)
	decrBytes := dongle.Decrypt.FromRawBytes(encrBytes).ByAes(cipher).ToBytes()
	paddRemovedBytes, err := RemovePKCS5Padding(decrBytes, AES_CBC_BLOCK)
	if err != nil {
		return nil, err
	}
	return paddRemovedBytes, nil
}

func Encrypt(plainText []byte, key []byte) (cypheredText []byte, err error) {
	stretchedKey := stretchKey(key)
	ivBytes := get16BytesIV()
	cipher := dongle.NewCipher()
	cipher.SetMode(dongle.CBC) // CBC、CFB、OFB、CTR、ECB
	//TODO: for some reason dongle PKCS5 padding is not working properly. Have to use custom padding function.
	//TODO: look into source code of dongle to see the problem
	cipher.SetPadding(dongle.No) // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cipher.SetKey(stretchedKey)  // key must be 16, 24 or 32 bytes
	cipher.SetIV(ivBytes)        // iv must be 16 bytes (ECB mode doesn't require setting iv)
	paddedText := PKCS5Padding(plainText, AES_CBC_BLOCK)
	encrBytes := dongle.Encrypt.FromBytes(paddedText).ByAes(cipher).ToRawBytes()
	allTheBytes := make([]byte, 0)
	allTheBytes = append(allTheBytes, ivBytes...)
	allTheBytes = append(allTheBytes, encrBytes...)
	return allTheBytes, nil
} //func encrypt(plainText []byte, key []byte) (cypheredText []byte, err error) {

func stretchKey(key []byte) []byte {
	const NO_OF_STRETCHING_ROUNDS = 128
	byteArray := dongle.Encrypt.FromBytes([]byte(key)).BySha512().ToRawBytes()
	for i := 1; i < 128; i++ {
		byteArray = dongle.Encrypt.FromBytes([]byte(byteArray)).BySha512().ToRawBytes()
	}
	return byteArray[:32]
} //func stretchKey(key string) []byte {

func RemovePKCS5Padding(paddedText []byte, blockSize int) ([]byte, error) {
	if len(paddedText) == 0 {
		return paddedText, nil
	}
	lastByte := paddedText[len(paddedText)-1]
	if int(lastByte) >= blockSize {
		return paddedText, nil
	}
	for i := 0; i < int(lastByte); i++ {
		penultimateByte := paddedText[len(paddedText)-1-i]
		if penultimateByte != lastByte {
			return paddedText, nil
		}
	}
	return paddedText[:len(paddedText)-int(lastByte)], nil
} //func RemovePKCS5Padding(paddedText []byte, blockSize int) ([]byte, error) {

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	// Code from here: https://gist.github.com/hothero/7d085573f5cb7cdb5801d7adcf66dcf3
	padding := blockSize - len(ciphertext)%blockSize
	if padding == blockSize {
		return ciphertext
	}
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// TODO: add error in return of this function
func get16BytesIV() []byte {
	//https://pkg.go.dev/crypto/rand#Read
	const IV_LENGTH = 16
	b := make([]byte, IV_LENGTH)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("get16BytesIV(). Error generating IV:%s\n", err)
		return nil
	}
	// The slice should now contain random bytes instead of only zeroes.
	// fmt.Println(bytes.Equal(b, make([]byte, c)))
	return b
} //func get16BytesIV() []byte {
