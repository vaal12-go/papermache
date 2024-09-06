package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

type MyWriteCloser struct {
	// MyWriteCloser is needed because go-qrcode/writer needs io.WriterCloser interface
	// https://stackoverflow.com/a/43115969
	io.Writer
}

func (mwc MyWriteCloser) Close() error {
	return nil
}

func QRCodeFromString(payload string) (*qrcode.QRCode, error) {
	//TODO: change error correction level depending on size of the payload
	qrc, err := qrcode.NewWith(payload,
		qrcode.WithEncodingMode(qrcode.EncModeAuto),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium))
	return qrc, err
}

func SendImage(w http.ResponseWriter, r *http.Request) {
	img_sha := r.URL.Path[10:]
	wc := MyWriteCloser{w}
	writer2 := standard.NewWithWriter(wc)
	cypherString, ok := imageMap[img_sha]
	if !ok {
		sendError(
			fmt.Sprintf("sendImage. Cannot find image with sha:%s\n", img_sha), &w)
	}
	qrc, err := QRCodeFromString(cypherString)
	if err != nil {
		sendError(
			fmt.Sprintf("sendImage. Error generating QRCode. err: %v\n", err), &w)
		return
	}
	if err = qrc.Save(writer2); err != nil {
		sendError(
			fmt.Sprintf("sendImage. Error writing QRCode image to ResponseWriter: %v\n", err), &w)
	}
} //func sendImage(w http.ResponseWriter, r *http.Request) {
