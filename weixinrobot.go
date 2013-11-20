package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

var token = "mumuxiaoxiaohai"

func main() {
	http.HandleFunc("/", search)
	log.Println("listen port 8044.")
	err := http.ListenAndServe(":8044", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func search(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	signature := r.Form.Get("signature")
	timestamp := r.Form.Get("timestamp")
	nonce := r.Form.Get("nonce")
	echostr := r.Form.Get("echostr")

	h := sha1.New()
	accessSlice := []string{token, timestamp, nonce}
	sort.Strings(accessSlice)
	io.WriteString(h, strings.Join(accessSlice, ""))

	sha1Str := hex.EncodeToString(h.Sum(nil))
	log.Println(signature, token, timestamp, nonce, sha1Str)

	// validate
	if sha1Str != signature {
		log.Println("validate failed.", echostr)
		os.Exit(1)
	}
	log.Println(r.Body)
}
