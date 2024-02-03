package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"runtime"
	"time"

	"github.com/golang-module/dongle"
)

//TODO: test openBrowser on Linux
//TODO: test openBrowser on macOs

// https://stackoverflow.com/questions/39320371/how-start-web-server-to-open-page-in-browser-in-golang
// open opens the specified URL in the default browser of the user.
func openBrowser(url string) error {
	time.Sleep(1 * time.Second)

	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows": //This works on Win10
		cmd = "explorer"
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, fmt.Sprintf("%s", url))
	return exec.Command(cmd, args...).Start()
} //func openBrowser(url string) error {

func encodeAndSendJSON(val any, w *http.ResponseWriter) {
	encrAnswerJSON, err := json.Marshal(val)
	// fmt.Printf("encrAnswerJSON: %v\n", string(encrAnswerJSON))
	if err != nil {
		io.WriteString(*w,
			fmt.Sprintf("Error converting string to JSON:%s", err))
	}
	io.WriteString(*w, string(encrAnswerJSON))
} //func encodeAndSendJSON(val any, w *http.ResponseWriter) {

func sendError(descr string, w *http.ResponseWriter) {
	fmt.Printf("Sending error to client:%s\n", descr)
	var encrAnswer = new(EncryptAnswer)
	encrAnswer.ErrOccurred = true
	encrAnswer.ErrDescription = descr
	encodeAndSendJSON(encrAnswer, w)
} //func sendError(descr string, w *http.ResponseWriter) {

func stretchKey(key string) []byte {
	const NO_OF_STRETCHING_ROUNDS = 128
	byteArray := dongle.Encrypt.FromBytes([]byte(key)).BySha512().ToRawBytes()

	for i := 1; i < 128; i++ {
		byteArray = dongle.Encrypt.FromBytes([]byte(byteArray)).BySha512().ToRawBytes()
	}
	// fmt.Printf("len of byteArray: %v\n", len(byteArray))
	// fmt.Printf("byteArray: %x\n", string(byteArray))
	return byteArray[:32]
} //func stretchKey(key string) []byte {

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

// Code from here: https://gist.github.com/hothero/7d085573f5cb7cdb5801d7adcf66dcf3
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// TODO: this is not good (replaced with stretchKey) - to be removed
// func padKey(key string) string {
// 	const KEY_PAD_SYMBOL = "="
// 	if len(key) < 16 {
// 		return key + strings.Repeat("=", 16-len(key))
// 	}
// 	return key
// }
