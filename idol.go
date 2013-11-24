package main

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	TypeNone        = 0
	TypeLegislation = 1 << (iota - 1)
	TypeCase
	TypeArticle
)

var ContentTypeMap = map[string]int{
	"搜法规": TypeLegislation,
	"搜案例": TypeCase,
	"搜评论": TypeArticle,
}

var ContentTypeDBMap = map[int]string{
	TypeLegislation: "law",
	TypeCase:        "case",
	TypeArticle:     "hotnews,ip_hottopic,ep_news_law,ep_news_case",
}

var ContentTypeTextMap = map[int]string{
	TypeLegislation: "法规",
	TypeCase:        "案例",
	TypeArticle:     "评论文章",
}

// 获取exp所对应的内容类型
func getContentType(exp string) int {
	if found, ok := ContentTypeMap[exp]; ok {
		return found
	}
	return TypeNone
}

func getAutnDatabaseName(t int) string {
	if found, ok := ContentTypeDBMap[t]; ok {
		return found
	}
	return "law"
}

func getContentTypeText(ct int) string {
	if found, ok := ContentTypeTextMap[ct]; ok {
		return found
	}
	return "法规"
}

func getArticles(keyword, db string) ([][]string, error) {
	autnResp := AutnResponse{}
	url := "http://192.168.2.210:9003/a=query&databasematch=%s" +
		"&sort=relevance+power_level:numberincreasing+date&print=none" +
		"&text=(%s):dretitle:drecontent:source:author:authorsource:articleid" +
		"+OR+(%s):tags"
	url = fmt.Sprintf(url, db, keyword, keyword)
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
		article := []string{hit.Title, "https://hk.lexiscn.com/" + hit.Reference}
		articles = append(articles, article)
	}

	return articles, nil
}
