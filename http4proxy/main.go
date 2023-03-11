package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/netip"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s source target\n", os.Args[0])
		return
	}

	source := os.Args[1]
	if !strings.Contains(source, ":") {
		source = "0.0.0.0:" + source
	} else if source[0] == ':' {
		source = "0.0.0.0" + source
	}

	addr, err := netip.ParseAddrPort(source)
	if err != nil {
		fmt.Println("invalid source", os.Args[1], err)
		return
	}

	target := os.Args[2]
	if !strings.Contains(target, ":") {
		target = "http://127.0.0.1:" + target
	} else if source[0] == ':' {
		target = "http://127.0.0.1" + target
	}

	targetURL, err := url.Parse(target)
	if err != nil {
		fmt.Println("invalid target", os.Args[2])
		return
	}

	fmt.Printf("reverse proxy: http://%v -> %v\n", addr, targetURL)
	err = http.ListenAndServe(addr.String(), httputil.NewSingleHostReverseProxy(targetURL))
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Println("closed")
	} else {
		fmt.Println(err)
	}
}
