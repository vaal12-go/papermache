package main

import (
	"embed"
	"errors"
	"fmt"
	_ "image/jpeg"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/vaal12-go/papermache/handlers"
)

//WinRes: https://github.com/tc-hib/go-winres

// https://pkg.go.dev/embed
var (
	//go:embed static
	content embed.FS
)

var srvr *http.Server

//[x]: create Shutdown button in html pages
//TODO: create option to create QR-Code from key as well
//TODO: create command line option for specific port
//TODO: create command line option for random port
//[x]: add reading of .env file for development

//LOW: to upgrade to 1.22 version for support of URLs in hanlers like: GET /

func createHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/stopServer", handlers.StopServer)
	mux.HandleFunc("/decodeQRCode", handlers.DecodeQRCode)
	mux.HandleFunc("/receiveData2Encrypt", handlers.ReceiveData2Encrypt)
	mux.HandleFunc("/qr-image/", handlers.SendImage)
	mux.HandleFunc("/index.html", handlers.RenderTemplate)
	mux.HandleFunc("/decode_qr.html", handlers.RenderTemplate)
	mux.HandleFunc("/decode_qr_result.html", handlers.RenderTemplate)
	mux.HandleFunc("/", handlers.GetRoot)
	return mux
}

// https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go
// URL above has an error, Handle Func should run with /static/
func main() {
	SetTitle("Papier-mâché 0.2.1 06Sep2024")
	err := godotenv.Load(".env")
	var DEVELOPER_MODE = ""
	//TODO: add -v parameter to show version of the executable
	var fs http.Handler = nil
	mux := createHandlers()
	if err != nil { //No .env file found - working as production
		//This is used for production
		fmt.Printf("%v\n", "Work in production mode.")
		fs = http.FileServer(http.FS(content))
		mux.Handle("/static/", fs)
		handlers.UseEmbeddedFS = true
		handlers.Content = content
	} else { //if err != nil {//.env file found
		DEVELOPER_MODE = os.Getenv("DEVELOPER_MODE")
		if DEVELOPER_MODE == "TRUE" { //DEVELOPER_MODE is set correctly - entering dev mode
			fmt.Printf("DEV mode\n")
			fs = http.FileServer(http.Dir("./static"))
			mux.Handle("/static/", http.StripPrefix("/static/", fs))
		} else { //Dev mode is set incorrectly - exiting
			fmt.Println(".env file is present, but no developer variable is set. Exisiting.")
		}
	} //} else { //if err != nil {//.env file found

	fmt.Println("Server starts.")
	go openBrowser("http://localhost:3333/index.html")
	srvr = &http.Server{
		Addr:    "localhost:3333",
		Handler: mux,
	}
	handlers.ServerInstance = srvr
	go terminalCharServerCloser(srvr)
	go func() {
		time.Sleep(time.Second * 2)
		fmt.Println("Press any key to stop server and exit.")
	}()
	err = srvr.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed gracefully. All is good.\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		if strings.Contains(err.Error(),
			"Only one usage of each socket address") {
			fmt.Println(PORT_ADDRESS_OCCUPIED_ERR_MSG)
		}
		fmt.Println("Press Enter key to exit")
		os.Exit(1)
	}
} //func main() {

//9cd7c8b6e0e7cf266accf920c4ec53d133e0e5f5eb3635506140ee8ef7a514d0
//34563fa7993331f673895a167bf2aab2542a409fe6bb765bf0226c582d54a7e226173c050e6f2c12d0d24b3a331fec0ae13617f328d2c192a4c3802fa21f06ae
