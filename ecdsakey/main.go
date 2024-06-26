package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"go.olapie.com/x/xconv"
	"go.olapie.com/x/xsecurity"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	pass := readConfirmedSecret("file password")
	if len(pass) < 8 {
		log.Println("Password is too short")
		return
	}
	pk := xconv.MustGet(xsecurity.GeneratePrivateKey(xsecurity.EcdsaP256))
	pri := xconv.MustGet(xsecurity.EncodePrivateKey(pk, pass))
	pub := xconv.MustGet(xsecurity.EncodePublicKey(xsecurity.GetPublicKey(pk)))
	name := time.Now().Format("20060102")
	_ = os.WriteFile(name+"-key.png", pri, 0644)
	_ = os.WriteFile(name+"-pub.png", pub, 0644)

	pubKey := xconv.MustGet(xsecurity.DecodePublicKey(pub)).(*ecdsa.PublicKey)
	priKey := xconv.MustGet(xsecurity.DecodePrivateKey(pri, pass)).(*ecdsa.PrivateKey)

	// Test
	hash := sha256.Sum256([]byte("message: hello"))
	sign := xconv.MustGet(ecdsa.SignASN1(rand.Reader, priKey, hash[:]))
	ok := ecdsa.VerifyASN1(pubKey, hash[:], sign)
	if !ok {
		log.Println("Test failed")
	}

	hash[0] = 20
	ok = ecdsa.VerifyASN1(pubKey, hash[:], sign)
	if ok {
		log.Println("Test failed")
	}
	log.Println("Test succeeded")
}

func readConfirmedSecret(name string) string {
	pass1 := readNonEmptyPassword(fmt.Sprintf("Enter %s: ", name))
	pass2 := readNonEmptyPassword(fmt.Sprintf("Repeat %s: ", name))
	if pass1 != pass2 {
		log.Println("Inputs mismatch")
		return ""
	}
	return pass1
}

func readNonEmptyPassword(msg ...any) string {
	var pass []byte
	for len(pass) == 0 {
		log.Print(msg...)
		pass = xconv.MustGet(terminal.ReadPassword(syscall.Stdin))
		log.Println()
	}
	return string(pass)
}
