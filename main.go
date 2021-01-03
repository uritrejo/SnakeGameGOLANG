package main

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"./controller"
	_ "./view"
)

const (
	port int = 4420
)

var (
	autoOpen bool = true
)

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)

	return exec.Command(cmd, args...).Start()
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	go controller.Start()

	fmt.Printf("Starting Snake webserver on port %d...\n", port)
	url := fmt.Sprintf("http://localhost:%d/", port)

	if autoOpen {
		fmt.Printf("Opening %s...\n", url)
		if err := open(url); err != nil {
			fmt.Println("Auto-open failed:", err)
			fmt.Printf("Open %s in your browser.\n", url)
		}
	} else {
		fmt.Printf("Auto-open not enabled, open %s in your browser.\n", url)
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

