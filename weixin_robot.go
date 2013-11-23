package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CmdFunc func(io.Writer, TextMessageReceived)

func cmdRoute(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// validate source
	if !validateSource(r) {
		return
	}

	// 读取信息内容
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Read message body error.", err)
		return
	}
	rcvMsg := BaseMessage{}
	xml.Unmarshal(body, &rcvMsg)
	fmt.Println("BaseMessage", rcvMsg)

	if rcvMsg.MsgType == "text" {
		textMsgRcv := TextMessageReceived{}
		xml.Unmarshal(body, &textMsgRcv)
		cmd := parseCmd([]byte(textMsgRcv.Content))

		textCmdRouter := map[int]CmdFunc{
			CmdSearch: search,
			CmdHelp:   help,
		}

		if f, ok := textCmdRouter[cmd]; ok {
			f(w, textMsgRcv)
			return
		}

		fmt.Println("cmd", cmd)
	} else if rcvMsg.MsgType == "event" {
		eventMsgRcv := SubscribeEventMessage{}
		xml.Unmarshal(body, &eventMsgRcv)

		fmt.Println("event message", eventMsgRcv)

		if eventMsgRcv.Event == "subscribe" {
			sendHelp(w, eventMsgRcv.FromUserName, eventMsgRcv.ToUserName)
		}
	}
}

func search(w io.Writer, rcvMsg TextMessageReceived) {
	msgContent := strings.Trim(rcvMsg.Content, " ")

	articles := [][]string{}

	sct := readContentType(rcvMsg.FromUserName)
	ct, keyword := parseSearchCmd([]byte(msgContent))

	fmt.Println("ct keyword", ct, keyword)

	db := "law"
	if ct != TypeNone {
		db = getAutnDatabaseName(ct)
	} else if sct != TypeNone {
		db = getAutnDatabaseName(sct)
	}
	if ct != sct {
		saveContentType(rcvMsg.FromUserName, ct)
	}

	log.Println("content_type", ct, sct, db, keyword)

	if "/:" != keyword {
		articles, _ = getTitles(keyword, db)
	}

	scopeText := getContentTypeText(ct)

	xmlb, err := genMsgContent(rcvMsg, articles, scopeText)
	if err != nil {
		log.Println("Generate xml error.", err)
		return
	}

	log.Println("xmlb", string(xmlb))
	fmt.Fprintf(w, string(xmlb))
}

func help(w io.Writer, rcvMsg TextMessageReceived) {
	sendHelp(w, rcvMsg.FromUserName, rcvMsg.ToUserName)
}

func main() {
	http.HandleFunc("/", cmdRoute)
	http.HandleFunc("/accountbindform", accountBindForm)
	log.Println("listen port 8044.")
	err := http.ListenAndServe(":8044", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
