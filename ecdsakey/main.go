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

	"go.olapie.com/security"
	"go.olapie.com/utils"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	pass := readConfirmedSecret("file password")
	if len(pass) < 8 {
		log.Println("Password is too short")
		return
	}
	pk := utils.MustGet(security.GeneratePrivateKey())
	pri := utils.MustGet(security.EncodePrivateKey(pk, pass))
	pub := utils.MustGet(security.EncodePublicKey(&pk.PublicKey))
	name := time.Now().Format("20060102")
	utils.MustNil(os.WriteFile(name+"-key.png", pri, 0644))
	utils.MustNil(os.WriteFile(name+"-pub.png", pub, 0644))

	pubKey := utils.MustGet(security.DecodePublicKey(pub))
	priKey := utils.MustGet(security.DecodePrivateKey(pri, pass))

	// Test
	hash := sha256.Sum256([]byte("message: hello"))
	sign := utils.MustGet(ecdsa.SignASN1(rand.Reader, priKey, hash[:]))
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
		pass = utils.MustGet(terminal.ReadPassword(syscall.Stdin))
		log.Println()
	}
	return string(pass)
}
