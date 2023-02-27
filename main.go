package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"golang.org/x/net/websocket"
)

func ShellServer(ws *websocket.Conn) {
	c := exec.Command("zsh")
	f, err := pty.Start(c)
	if err != nil {
		ws.Write([]byte(fmt.Sprintf("Error creating pty: %s\r\n", err)))
		ws.Close()
		return
	}

	go io.Copy(ws, f)
	io.Copy(f, ws)
	ws.Close()
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./")))
	http.Handle("/echo", websocket.Handler(ShellServer))
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		panic(err)
	}
	os.Exit(0)
}
