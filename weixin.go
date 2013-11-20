package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const DataDir = "data/"

func genMsgContent(msgRcv TextMessageReceived, articles [][]string, scope string) ([]byte, error) {
	articleCount := len(articles) + 1
	newsMsgItems := make([]NewsMessageItem, articleCount)
	newsMsgItems[0] = NewsMessageItem{
		xml.Name{"", "item"}, "范围:" + scope + "  搜索结果", "", "", "",
	}
	if articleCount > 1 {
		for i, article := range articles {
			newsMsgItems[i+1] = NewsMessageItem{
				xml.Name{"", "item"}, article[0], "", "",
				"http://www.lexiscn.com/" + strings.Trim(article[1], "/")}
		}
	} else {
		tmpNewsMsgItem := NewsMessageItem{
			xml.Name{"", "item"}, "无结果", "", "", "",
		}
		newsMsgItems = append(newsMsgItems, tmpNewsMsgItem)
		articleCount += 1
	}
	newsMsgArticle := NewsMessageArtice{
		xml.Name{"", "Articles"}, newsMsgItems,
	}
	newsMsg := NewsMessage{
		xml.Name{"", "xml"},
		msgRcv.FromUserName,
		msgRcv.ToUserName,
		time.Now().Unix(),
		"news",
		articleCount,
		newsMsgArticle,
	}
	b, err := xml.Marshal(newsMsg)
	if err != nil {
		return []byte(""), err
	}
	return b, nil
}

func validateSource(r *http.Request) bool {
	r.ParseForm()
	signature := r.Form.Get("signature")
	timestamp := r.Form.Get("timestamp")
	nonce := r.Form.Get("nonce")
	echostr := r.Form.Get("echostr")

	h := sha1.New()
	accessSlice := []string{Token, timestamp, nonce}
	sort.Strings(accessSlice)
	io.WriteString(h, strings.Join(accessSlice, ""))
	sha1Str := hex.EncodeToString(h.Sum(nil))

	log.Println(signature, Token, timestamp, nonce, sha1Str)

	// validate
	if sha1Str != signature {
		log.Println("validate failed.", echostr)
		return false
	}
	return true
}

func saveContentType(user string, contentType int) bool {
	c := fmt.Sprintf("content_type %d", contentType)

	err := ioutil.WriteFile(DataDir+user, []byte(c), os.ModePerm)
	if err != nil {
		log.Println("write file failed.")
	}
	return true
}

func readContentType(user string) int {
	ct := TypeNone
	b, err := ioutil.ReadFile(DataDir + user)
	if err != nil {
		return ct
	}
	s := strings.Fields(string(b))
	log.Println("content_from_file", s)
	if len(s) < 2 {
		return ct
	}
	ct, err = strconv.Atoi(s[1])
	if err != nil {
		return ct
	}
	return ct
}
