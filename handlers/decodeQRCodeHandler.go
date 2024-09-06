package handlers

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"net/http"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/vaal12-go/papermache/encrypt"
)

func decodeQRImage(imageBytes []byte) (*gozxing.Result, error) {
	bReader := bytes.NewReader(imageBytes)
	img, _, err := image.Decode(bReader)
	if err != nil {
		return nil, err
	}
	bmp, _ := gozxing.NewBinaryBitmapFromImage(img)
	qrReader := qrcode.NewQRCodeReader()
	result, err := qrReader.Decode(bmp, map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER:   true,
		gozxing.DecodeHintType_PURE_BARCODE: true,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DecodeQRCode(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10485760) //10MB file in memory (this does not limit upload)
	file, _, err := r.FormFile("qr_file")
	if err != nil {
		fmt.Printf("file reading err: %v\n", err)
		sendError(
			fmt.Sprintf("Error file in request: %v", err), &w)
	}
	defer file.Close()
	var buf bytes.Buffer
	io.Copy(&buf, file)
	keyString := r.FormValue("key")
	result, err := decodeQRImage(buf.Bytes())
	if err != nil {
		fmt.Printf("qrReader.Decode err: %v\n", err)
		sendError(
			fmt.Sprintf("Error decoding image in request err: %v", err), &w)
		return
	}
	// else {
	// 	fmt.Println("QR code decoded:" + result.String())
	// }

	if keyString == "" {
		ans := map[string]string{
			"answer": result.GetText(),
		}
		encodeAndSendJSON(ans, &w)
	} else { //if keyString == "" {
		secrStr, err := encrypt.DecodeCode(result.GetText(), keyString)
		if err != nil {
			ans := map[string]string{
				"answer": "",
				"error":  err.Error(),
			}
			encodeAndSendJSON(ans, &w)
			return
		}
		ans := map[string]string{
			"answer": secrStr,
		}
		encodeAndSendJSON(ans, &w)
	} //} else { //if keyString == "" {
} //func decodeQRCode(w http.ResponseWriter, r *http.Request) {
