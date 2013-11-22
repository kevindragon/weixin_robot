package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var Token = "mumuxiaoxiaohai"

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
	echostr := r.Form.Get("echostr")

	// validate source
	if !validateSource(r) {
		log.Println("validate failed.", echostr)
		return
	}

	// 读取信息内容
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Read body", err)
		return
	}
	msgRcv := TextMessageReceived{}
	xml.Unmarshal(body, &msgRcv)

	log.Println("message", msgRcv)

	if msgRcv.MsgType == "text" {
		msgContent := strings.Trim(msgRcv.Content, " ")
		if msgContent == "帮助" || msgContent == "?" || msgContent == "？" {
			helpMsgXml, err := genTextMsgContent(msgRcv.ToUserName, msgRcv.FromUserName, helpMessage())
			if err != nil {
				log.Println("generate help message xml error.", err)
			}
			log.Println("helpMsgXml", string(helpMsgXml))
			fmt.Fprintf(w, string(helpMsgXml))
			return
		}

		articles := [][]string{}

		keyword := msgContent
		keywordRune := []rune(keyword)

		sct := readContentType(msgRcv.FromUserName)
		var ct int
		if len(keywordRune) >= 3 {
			ct = getContentType(string(keywordRune[:3]))
		} else {
			ct = getContentType(string(keywordRune))
		}
		var db string
		if ct == TypeNone {
			if sct == TypeNone {
				db = "law"
			} else {
				ct = sct
				db = getAutnDatabaseName(sct)
			}
		} else {
			db = getAutnDatabaseName(ct)
			if len(keywordRune) > 2 {
				if len(keywordRune) >= 3 {
					keyword = strings.Trim(string(keywordRune[3:]), " ")
				} else {
					keyword = strings.Trim(string(keywordRune), "")
				}
			}
			if ct != sct {
				saveContentType(msgRcv.FromUserName, ct)
			}
		}

		log.Println("content_type", ct, sct, db, keyword)

		if "/:" != keyword {
			articles, _ = getTitles(keyword, db)
		}

		scope := getContentTypeText(ct)

		xmlb, err := genMsgContent(msgRcv, articles, scope)
		if err != nil {
			log.Println("Generate xml error.", err)
			return
		}

		log.Println("xmlb", string(xmlb))
		fmt.Fprintf(w, string(xmlb))
	} else if msgRcv.MsgType == "event" {
		eventMsgRcv := SubscribeEventMessage{}
		xml.Unmarshal(body, &eventMsgRcv)
		log.Println(eventMsgRcv)
		if eventMsgRcv.Event == "subscribe" {
			sendHelp(w, eventMsgRcv.FromUserName, eventMsgRcv.ToUserName)
		}
	}
}