package main

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"syscall"
	"unsafe"
)

//https://pkg.go.dev/embed
//go:embed static/**

var content embed.FS

//TODO: create Shutdown button in html pages
//TODO: create option to create QR-Code from key as well
//TODO: create command line option for specific port
//TODO: create command line option for random port
//TODO: add reading of .env file for development

//WinRes: https://github.com/tc-hib/go-winres

// From here: https://github.com/lxi1400/GoTitle/blob/main/title.go
// https://www.reddit.com/r/golang/comments/a51266/how_get_or_set_the_console_title_in_go/
// TODO: add check if this is windows maching and find solution for linux
func SetTitle(title string) (int, error) {
	handle, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, err
	}
	defer syscall.FreeLibrary(handle)
	proc, err := syscall.GetProcAddress(handle, "SetConsoleTitleW")
	if err != nil {
		return 0, err
	}
	r, _, err := syscall.Syscall(proc, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	return int(r), err
}

// https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go
// URL above has an error, Handle Func should run with /static/
func main() {

	SetTitle("Papier-mache 0.1.1")

	mux := http.NewServeMux()

	//This is used for development
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	//TODO: review and fix any errors - some files are not seen
	//This is used for production
	// fs := http.FileServer(http.FS(content))
	// mux.Handle("/static/", fs)

	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/receiveData2Encrypt", receiveData2Encrypt)
	mux.HandleFunc("/qr-image/", sendImage)

	url_2open := "http://localhost:3333/static/index.html"
	fmt.Printf("Starting server. Please point your server to:%s", url_2open)
	go openBrowser(url_2open)

	err := http.ListenAndServe("localhost:3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed gracefully. All is good.\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
} //func main() {
