package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"
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
	echostr := r.Form.Get("echostr")

	// validate source
	if !validateSource(r) {
		log.Println("validate failed.", echostr)
		return
	}

	// 读取信息内容
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Read body", err)
		return
	}
	msgRcv := TextMessageReceived{}
	xml.Unmarshal(b, &msgRcv)

	if msgRcv.MsgType == "text" {
		log.Println("content", msgRcv.Content)
		articles := [][]string{}
		if "/:" != msgRcv.Content[:2] {
			articles, _ = getTitles(msgRcv.Content)
		}

		xmlb, err := genMsgContent(msgRcv, articles)
		if err != nil {
			log.Println("Generate xml error.", err)
			return
		}

		log.Println("xmlb", string(xmlb))
		fmt.Fprintf(w, string(xmlb))
	}
}

func genMsgContent(msgRcv TextMessageReceived, articles [][]string) ([]byte, error) {
	articleCount := len(articles) + 1
	newsMsgItems := make([]NewsMessageItem, articleCount)
	newsMsgItems[0] = NewsMessageItem{xml.Name{"", "item"}, "搜索结果", "", "", ""}
	if articleCount > 1 {
		for i, article := range articles {
			newsMsgItems[i+1] = NewsMessageItem{
				xml.Name{"", "item"},
				article[0],
				"",
				"",
				"http://www.lexiscn.com/" + strings.Trim(article[1], "/")}
		}
	} else {
		newsMsgItems = append(newsMsgItems, NewsMessageItem{xml.Name{"", "item"}, "无结果", "", "", ""})
		articleCount += 1
	}
	newsMsgArticle := NewsMessageArtice{
		xml.Name{"", "Articles"},
		newsMsgItems,
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

func getTitles(keyword string) ([][]string, error) {
	autnResp := AutnResponse{}
	url := "http://192.168.2.210:9003/a=query&databasematch=law" +
		"&sort=relevance+power_level:numberincreasing+date&print=none" +
		"&text=(%s):dretitle:articleid+OR+(%s):tags"
	url = fmt.Sprintf(url, keyword, keyword)
	log.Println("url", url)
	resp, err := http.Get(url)
	if err != nil {
		return [][]string{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return [][]string{}, err
	}

	xml.Unmarshal(body, &autnResp)

	articles := [][]string{}
	if autnResp.Response != "SUCCESS" {
		log.Println("Search engine response error.")
		return articles, errors.New("Response failed.")
	}

	for _, hit := range autnResp.Respdata.Hits {
		article := []string{hit.Title, hit.Reference}
		articles = append(articles, article)
	}

	return articles, nil
}

func validateSource(r *http.Request) bool {
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
		return false
	}
	return true
}
