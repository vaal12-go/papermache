package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

// This is needed because go-qrcode/writer needs io.WriterCloser interface
// https://stackoverflow.com/a/43115969
type MyWriteCloser struct {
	io.Writer
}

func (mwc MyWriteCloser) Close() error {
	// Noop
	return nil
}

func sendImage(w http.ResponseWriter, r *http.Request) {
	// fmt.Printf("URL:%s\n", r.URL.Path)
	img_sha := r.URL.Path[10:]
	// fmt.Printf("img_sha: %v\n", img_sha)
	// io.WriteString(w, "hello I am a qrimage:"+img_sha+
	// 	"\n"+imageMap[img_sha])

	// var w io.WriterCloser
	wc := MyWriteCloser{w}
	writer2 := standard.NewWithWriter(wc)

	cypherText, ok := imageMap[img_sha]
	if !ok {
		sendError(
			fmt.Sprintf("sendImage. Cannot find image with sha:%s\n", img_sha), &w)
	}

	qrc, err := qrcode.New(cypherText)
	if err != nil {
		sendError(
			fmt.Sprintf("sendImage. Error generating QRCode. err: %v\n", err), &w)
		return
	}

	if err = qrc.Save(writer2); err != nil {
		// fmt.Printf("could not save image: %v", err)
		sendError(
			fmt.Sprintf("sendImage. Error writing QRCode image to ResponseWriter: %v\n", err), &w)
	}
} //func sendImage(w http.ResponseWriter, r *http.Request) {
