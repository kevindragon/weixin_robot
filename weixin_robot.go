package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CmdFunc func(io.Writer, TextMessageReceived)

func cmdRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Println("r.URL", r.URL)
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
			CmdSearch:      search,
			CmdHelp:        help,
			CmdAccountBind: sendAccountBindLink,
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

	sct := readContentType(rcvMsg.FromUserName)
	ct, keyword := parseSearchCmd([]byte(msgContent))
	if ct != TypeNone && ct != sct {
		saveContentType(rcvMsg.FromUserName, ct)
	}

	db := "law"
	fmt.Println(ct, TypeNone, ct == TypeNone, sct, sct == ct)
	if ct == TypeNone && sct == TypeNone {
		ct = TypeLegislation
	} else if ct != TypeNone {
		db = getAutnDatabaseName(ct)
	} else if sct != TypeNone {
		db = getAutnDatabaseName(sct)
		ct = sct
	}

	log.Println("content_type", ct, sct, db, keyword)

	if 0 != strings.Index(msgContent, "/:") {
		articles, _ := getArticles(keyword, db)

		scopeText := getContentTypeText(ct)
		items := [][]string{[]string{"范围:" + scopeText + "  搜索结果", ""}}
		if len(articles) > 0 {
			for _, article := range articles {
				items = append(items, article)
			}
		} else {
			items = append(items, []string{"无结果", ""})
		}

		xmlb, err := genTeleTextMsgContent(rcvMsg, items)
		if err != nil {
			log.Println("Generate xml error.", err)
			return
		}

		log.Println("xmlb", string(xmlb))
		fmt.Fprintf(w, string(xmlb))
	} else {
		fmt.Println("发送的是表情")
	}
}

func sendAccountBindLink(w io.Writer, rcvMsg TextMessageReceived) {
	items := [][]string{[]string{"绑定账号", "http://staging2.lexisnexis.com.cn/weixin/accountbindform"}}
	fmt.Println(items)
	xmlb, err := genTeleTextMsgContent(rcvMsg, items)
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

func accountBindForm(w http.ResponseWriter, r *http.Request) {
	fmt.Println("in accountBindForm")
	t, err := template.ParseFiles("templates/accountbindform.html")
	if err != nil {
		log.Println("parse files error.", err)
		return
	}
	t.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", cmdRoute)
	http.HandleFunc("/accountbindform", accountBindForm)
	http.Handle("/favicon.ico", http.FileServer(http.Dir("./")))
	log.Println("listen port 8044.")
	err := http.ListenAndServe(":8044", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
