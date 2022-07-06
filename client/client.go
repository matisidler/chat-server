package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

var (
	port = flag.Int("p", 3090, "port to listen on")
	host = flag.String("h", "localhost", "host to listen on")
)

func main() {
	flag.Parse()
	// Listen on TCP port
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		panic(err)
	}
	fmt.Println("connected")
	done := make(chan struct{})
	// Start a goroutine to receive data from the server
	go func() {
		io.Copy(os.Stdout, conn)
		done <- struct{}{}
	}()
	// Copy what we got to the console line
	CopyContent(conn, os.Stdin)
	conn.Close()
	<-done
}

// CopyCOntent copies content from src to dst
func CopyContent(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		fmt.Fprintf(os.Stderr, "io.Copy: %v\n", err)
		os.Exit(1)
	}
}
