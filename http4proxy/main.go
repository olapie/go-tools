package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"go.olapie.com/tools/http4proxy/asset"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

var arguments = struct {
	certFile     string
	keyFile      string
	port         int
	securityPort int
}{}

func main() {
	flag.StringVar(&arguments.certFile, "cert", "", "path to certificate pem")
	flag.StringVar(&arguments.keyFile, "key", "", "ath to key pem")
	flag.IntVar(&arguments.port, "port", 8080, "Proxy server port")
	flag.IntVar(&arguments.securityPort, "security_port", 8443, "Proxy server security port")
	flag.Parse()

	// 没有作用，客户端发送的非TLS形式的CONNECT命令，该代理就无法处理
	//go startHTTPSProxy()
	startHTTPProxy()
}

func startHTTPProxy() {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", arguments.port),
		Handler: http.HandlerFunc(handle),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Println("Proxy server listening at", server.Addr)
	log.Fatal(server.ListenAndServe())
}

func startHTTPSProxy() {
	if arguments.certFile != "" {
		server := &http.Server{
			Addr:    fmt.Sprintf(":%d", arguments.securityPort),
			Handler: http.HandlerFunc(handle),
			//// Disable HTTP/2.
			//TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
			//TLSConfig: &tls.Config{
			//	MinVersion: tls.VersionTLS13,
			//},
		}
		log.Println(arguments.certFile, arguments.keyFile)
		log.Println("Proxy secure server listening at", server.Addr)
		log.Fatal(server.ListenAndServeTLS(arguments.certFile, arguments.keyFile))
		return
	}

	cert, err := tls.X509KeyPair(asset.CertPEM, asset.KeyPEM)
	if err != nil {
		log.Fatal(err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS13,
	}

	server := http.Server{
		TLSConfig: tlsConfig,
		Addr:      fmt.Sprintf(":%d", arguments.securityPort),
		Handler:   http.HandlerFunc(handle),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}
	log.Println("Proxy secure server listening at", server.Addr)
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func handle(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.URL)
	if r.Method == http.MethodConnect {
		handleTunneling(w, r)
	} else {
		handleHTTP(w, r)
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	destConn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(destConn, clientConn)
	go transfer(clientConn, destConn)
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
