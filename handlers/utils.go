package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func saveQRCodeToByteArray(code *qrcode.QRCode) ([]byte, error) {
	imageBuf := bytes.Buffer{}
	wc := MyWriteCloser{&imageBuf}
	stdWriter := standard.NewWithWriter(wc)
	if err := code.Save(stdWriter); err != nil {
		return nil, err
	}
	return imageBuf.Bytes(), nil
}

func encodeAndSendJSON(val any, w *http.ResponseWriter) {
	encrAnswerJSON, err := json.Marshal(val)
	if err != nil {
		io.WriteString(*w,
			fmt.Sprintf("Error converting string to JSON:%s", err))
	}
	io.WriteString(*w, string(encrAnswerJSON))
} //func encodeAndSendJSON(val any, w *http.ResponseWriter) {

func sendError(descr string, w *http.ResponseWriter) {
	// fmt.Printf("Sending error to client:%s\n", descr)
	var encrAnswer = new(EncryptAnswer)
	encrAnswer.ErrOccurred = true
	encrAnswer.ErrDescription = descr
	encodeAndSendJSON(encrAnswer, w)
} //func sendError(descr string, w *http.ResponseWriter) {
