package main

import (
	"log"
	"os"

	"github.com/things-go/go-socks5"
)

/**
  SOCKS stands for SocketSecure
  SOCKS5 supports TCP/UDP traffic, while SOCKS4 only supports TCP
  Refer to https://oxylabs.io/blog/socks-vs-http-proxy
*/

func main() {
	// Create a SOCKS5 server
	server := socks5.NewServer(
		socks5.WithLogger(socks5.NewLogger(log.New(os.Stdout, "socks5: ", log.LstdFlags))),
	)

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", ":8000"); err != nil {
		panic(err)
	}
}
