package main

import (
	"flag"
	"fmt"
	"go.olapie.com/utils"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"go.olapie.com/ola/headers"
)

func main() {
	var addr string
	var dir string
	localIP := utils.GetPrivateIPv4(utils.GetPrivateIPv4Interface())
	flag.StringVar(&addr, "addr", localIP.String()+":0", "address")
	flag.StringVar(&dir, "dir", ".", "directory")
	flag.Parse()

	user := fmt.Sprint(rand.Int31() % 1e3)
	password := fmt.Sprint(rand.Int31() % 1e3)
	fmt.Println("User:", user)
	fmt.Println("Password:", password)
	userToAuthorization := headers.CreateUserAuthorizations(map[string]string{user: password})

	fs := http.FileServer(http.Dir(dir))
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("http://" + l.Addr().String())
	err = http.Serve(l, http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		reqAuth := headers.GetAuthorization(req.Header)
		for u, auth := range userToAuthorization {
			if auth == reqAuth {
				log.Println(u)
				fs.ServeHTTP(w, req)
				return
			}
		}
		w.Header().Set(headers.KeyWWWAuthenticate, "Basic realm="+strconv.Quote("olapie"))
		w.WriteHeader(http.StatusUnauthorized)
	}))
	log.Println(err)
}
