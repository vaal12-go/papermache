package main

import (
	"bufio"
	"embed"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/joho/godotenv"
)

//https://pkg.go.dev/embed
//go:embed static/**

var content embed.FS

const PORT_ADDRESS_OCCUPIED_ERR_MSG = `
This error is usually due to other copy of this programm running.
Either exit other copy and start this program once more or just point your browser to
	http://localhost:3333/
if other copy works correctly you should be able to use it.

To exit other program press Ctrl-C in console window where it is running.
`

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

func readSingleCharFromConsole() {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()

	if err != nil {
		fmt.Println(err)
	}

	// print out the unicode value i.e. A -> 65, a -> 97
	fmt.Println(char)

	switch char {
	case 'A':
		fmt.Println("A Key Pressed")
		break
	case 'a':
		fmt.Println("a Key Pressed")
		break
	}
}

// https://www.digitalocean.com/community/tutorials/how-to-make-an-http-server-in-go
// URL above has an error, Handle Func should run with /static/
func main() {

	SetTitle("Papier-mâché 0.1.2")

	mux := http.NewServeMux()

	err := godotenv.Load(".env")
	// dbgMode := false
	var DEVELOPER_MODE = ""

	var fs http.Handler = nil

	if err != nil { //No .env file found - working as production
		//TODO: review and fix any errors - some files are not seen
		//This is used for production
		fmt.Printf("%v\n", "Work in production mode.")
		fs = http.FileServer(http.FS(content))
		mux.Handle("/static/", fs)

	} else { //if err != nil {//.env file found
		DEVELOPER_MODE = os.Getenv("DEVELOPER_MODE")
		if DEVELOPER_MODE == "TRUE" { //DEVELOPER_MODE is set correctly - entering dev mode
			fmt.Printf("DEV mode\n")
			//This is used for development
			fs = http.FileServer(http.Dir("./static"))
			mux.Handle("/static/", http.StripPrefix("/static/", fs))
		} else { //Dev mode is set incorrectly - exiting
			fmt.Println(".env file is present, but no developer variable is set. Exisiting.")
		}
	} //} else { //if err != nil {//.env file found

	mux.HandleFunc("/", getRoot)
	mux.HandleFunc("/receiveData2Encrypt", receiveData2Encrypt)
	mux.HandleFunc("/qr-image/", sendImage)

	fmt.Printf("\"Starting server\": %v\n", "Starting server")

	go openBrowser("http://localhost:3333/static/index.html")

	err = http.ListenAndServe("localhost:3333", mux)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Server closed gracefully. All is good.\n")
	} else if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		if strings.Contains(err.Error(), "Only one usage of each socket address") {
			fmt.Println(PORT_ADDRESS_OCCUPIED_ERR_MSG)
		}

		fmt.Println("Press Enter key to exit")
		readSingleCharFromConsole()
		os.Exit(1)
	}
} //func main() {

//SHA256 - 9cd7c8b6e0e7cf266accf920c4ec53d133e0e5f5eb3635506140ee8ef7a514d0
//Service
