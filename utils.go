package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/term"
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

func terminalCharServerCloser(srvr *http.Server) {
	readSingleCharFromConsole()
	srvr.Shutdown(context.TODO())
}

func readSingleCharFromConsole() {
	// fmt.Printf("\"waiting for rune\": %v\n", "waiting for rune")
	// reader := bufio.NewReader(os.Stdin)
	// char, _, err := reader.ReadRune()

	// if err != nil {
	// 	fmt.Println(err)
	// }
	// // print out the unicode value i.e. A -> 65, a -> 97
	// fmt.Printf("Have rune:%v:", char)
	// switch char {
	// case 'A':
	// 	fmt.Println("A Key Pressed")
	// 	break
	// case 'a':
	// 	fmt.Println("a Key Pressed")
	// 	break
	// }
	// deadlinetime := time.Now().Add(time.Second * 5)
	// // err := os.Stdin.SetDeadline(deadlinetime)

	// // if err != nil {
	// // 	fmt.Printf("Setting deadline err: %v\n", err)
	// // 	return
	// // }

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	b := make([]byte, 1)
	_, err = os.Stdin.Read(b)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("the char %q was hit", string(b[0]))
}

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
