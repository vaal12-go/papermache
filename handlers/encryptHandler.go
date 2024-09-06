package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-module/dongle"
	"github.com/vaal12-go/papermache/encrypt"
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

func verifyQRCode(qrCode *qrcode.QRCode, encodedString string) (verified bool, err error) {
	imgBuf, err := saveQRCodeToByteArray(qrCode)
	if err != nil {
		return false, fmt.Errorf("Error saving QRCode as stream. Err:%v\n", err)
	}

	result, err := decodeQRImage(imgBuf)

	//[x]: move below to separate function (reading of QR code from byte array)
	// bReader := bytes.NewReader(imgBuf)
	// img, _, err := image.Decode(bReader)
	// if err != nil {
	// 	return false, fmt.Errorf("Error decoding image in request err:%v\n", err)
	// }
	// bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	// qrReader := zxingqrcode.NewQRCodeReader()
	// result, err := qrReader.Decode(bmp, map[gozxing.DecodeHintType]interface{}{
	// 	gozxing.DecodeHintType_TRY_HARDER:   true,
	// 	gozxing.DecodeHintType_PURE_BARCODE: true,
	// })
	if err != nil {
		return false, fmt.Errorf("Error decoding image Code:%v\n", err)
	}
	if strings.Compare(result.String(), encodedString) == 0 {
		// fmt.Printf("\"QR generation succeeded\": %v\n", "QR generation succeeded")
		return true, nil
	} else {
		return false, nil
	}
} //func verifyQRCode(qrCode *qrcode.QRCode, encodedString string) (verified bool, err error) {

func ReceiveData2Encrypt(w http.ResponseWriter, r *http.Request) {
	var bb bytes.Buffer
	bb.ReadFrom(r.Body)
	var p RequestToEncrypt
	err := json.Unmarshal(bb.Bytes(), &p)
	if err != nil {
		sendError(
			fmt.Sprintf("receiveData2Encrypt. Error unmarshalling JSON. err: %v\n", err), &w)
		return
	}
	// fmt.Printf("p.Text2Encrypt: %v\n", p.Text2Encrypt)
	allTheBytes, err := encrypt.Encrypt([]byte(p.Text2Encrypt), []byte(p.Key))
	if err != nil {
		sendError(
			fmt.Sprintf("receiveData2Encrypt. Error encrypting. err: %v\n", err), &w)
		return
	}
	encryptedBase64String := dongle.Encode.FromBytes(allTheBytes).ByBase64().ToString()
	qrCode, err := QRCodeFromString(encryptedBase64String)
	if err != nil {
		sendError(
			fmt.Sprintf("receiveData2Encrypt. Error generating QRCode. err: %v\n", err), &w)
		return
	}
	//Check QRCode can be read back
	res, err := verifyQRCode(qrCode, encryptedBase64String)
	if err != nil {
		sendError(
			fmt.Sprintf("Error checking QR CODE. err: %v\n", err), &w)
		return
	}
	if !res {
		sendError(
			fmt.Sprintf("QRCode does not verify. err: %v\n", err), &w)
		return
	}
	//TODO: generate image name based on sha-512
	img_sha := dongle.Encrypt.FromBytes([]byte(allTheBytes)).BySha512().ToHexString()
	for k := range imageMap { //Clearing imageMap so it will not leak memory
		delete(imageMap, k)
	}
	imageMap[img_sha] = encryptedBase64String

	var encrAnswer = new(EncryptAnswer)
	encrAnswer.EncryptedText = encryptedBase64String
	encrAnswer.ImagePath = "/qr-image/" + img_sha
	encodeAndSendJSON(encrAnswer, &w)
} //func receiveData2Encrypt(w http.ResponseWriter, r *http.Request) {
