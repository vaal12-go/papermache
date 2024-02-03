package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-module/dongle"
	"github.com/yeqown/go-qrcode/v2"
)

type RequestToEncrypt struct {
	Text2Encrypt string
	Key          string
}

type EncryptAnswer struct {
	EncryptedText  string
	ImagePath      string
	ErrOccurred    bool
	ErrDescription string
}

var imageMap = make(map[string]string, 0) //This struct will share Cyphers between encryptions so sendImage maywork

func receiveData2Encrypt(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("got /receiveData2Encrypt request\n")
	var bb bytes.Buffer
	bb.ReadFrom(r.Body)
	var p RequestToEncrypt
	err := json.Unmarshal(bb.Bytes(), &p)
	if err != nil {
		sendError(
			fmt.Sprintf("receiveData2Encrypt. Error unmarshalling JSON. err: %v\n", err), &w)
		return
	}

	// fmt.Printf("p: %#v\n", p)
	// paddedKey := padKey(p.Key)
	// fmt.Printf("paddedKey: %v\n", paddedKey)
	stretchedKey := stretchKey(p.Key)
	// fmt.Printf("stretchedKey: %v\n", stretchedKey)

	ivBytes := get16BytesIV()
	// fmt.Printf("ivBytes: %v\n", ivBytes)

	cipher := dongle.NewCipher()
	cipher.SetMode(dongle.CBC) // CBC、CFB、OFB、CTR、ECB

	//TODO: for some reason dongle PKCS5 padding is not working properly. Have to use custom padding function.
	//TODO: look into source code of dongle to see the problem
	cipher.SetPadding(dongle.No) // No、Empty、Zero、PKCS5、PKCS7、AnsiX923、ISO97971
	cipher.SetKey(stretchedKey)  // key must be 16, 24 or 32 bytes
	cipher.SetIV(ivBytes)        // iv must be 16 bytes (ECB mode doesn't require setting iv)

	// fmt.Printf("\"Before encryption\": %v\n", "Before encryption")
	// fmt.Printf("p.Text2Encrypt length: %v\n", len(p.Text2Encrypt))

	paddedText := PKCS5Padding([]byte(p.Text2Encrypt), 16)

	// fmt.Printf("paddedText length: %v\n", len(paddedText))
	encrBytes := dongle.Encrypt.FromBytes(paddedText).ByAes(cipher).ToRawBytes()
	// encrBytes := dongle.Encrypt.FromString(p.Text2Encrypt).ByAes(cipher).ToRawBytes()
	// fmt.Printf("encrBytes length: %v\n", len(encrBytes))

	// ivBytesBase64String := dongle.Encode.FromBytes(ivBytes).ByBase64().ToString()

	allTheBytes := make([]byte, 0)

	allTheBytes = append(allTheBytes, ivBytes...)
	allTheBytes = append(allTheBytes, encrBytes...)
	// fmt.Printf("allTheBytes: %v\n", allTheBytes)
	// fmt.Printf("allTheBytes len: %v\n", len(allTheBytes))

	// Encrypt by aes from string and output raw string
	encryptedBase64String := dongle.Encode.FromBytes(allTheBytes).ByBase64().ToString()
	// fmt.Printf("encryptedBase64String length: %v\n", len(encryptedBase64String))

	// qrc, err := qrcode.New(encryptedBase64String)
	_, err = qrcode.New(encryptedBase64String)
	if err != nil {
		// fmt.Printf("Cannot generate QRCode: %v\n", err)
		sendError(
			fmt.Sprintf("receiveData2Encrypt. Error generating QRCode. err: %v\n", err), &w)
		return
	}

	//TODO: generate image name based on sha-512
	//Below is not needed due image generation on the fly
	// writ, err = standard.New("./static/qr-codes/repo-qrcode.jpeg")
	// if err != nil {
	// 	fmt.Printf("standard.New failed: %v", err)
	// 	sendError(
	// 		fmt.Sprintf("Error creating writer for QRCode. err: %v\n", err), &w)
	// 	return
	// }

	// save file
	//Only needed for debug - image is generated on the fly by sendImage handler
	// if err = qrc.Save(writ); err != nil {
	// 	// fmt.Printf("could not save image: %v", err)
	// 	sendError(
	// 		fmt.Sprintf("Error saving QRCode image. err: %v\n", err), &w)
	// }

	var encrAnswer = new(EncryptAnswer)

	img_sha := dongle.Encrypt.FromBytes([]byte(allTheBytes)).BySha512().ToHexString()
	// fmt.Printf("img_sha: %v\n", img_sha)

	for k := range imageMap { //Clearing imageMap so it will not leak memory
		delete(imageMap, k)
	}
	imageMap[img_sha] = encryptedBase64String
	encrAnswer.EncryptedText = encryptedBase64String
	encrAnswer.ImagePath = "/qr-image/" + img_sha
	encodeAndSendJSON(encrAnswer, &w)
} //func receiveData2Encrypt(w http.ResponseWriter, r *http.Request) {
