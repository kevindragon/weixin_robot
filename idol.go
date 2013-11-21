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

func getContentType(exp string) int {
	m := map[string]int{
		"搜法规": TypeLegislation,
		"搜案例": TypeCase,
		"搜评论": TypeArticle,
	}
	if found, ok := m[exp]; ok {
		return found
	}
	return TypeNone
}

func getAutnDatabaseName(t int) string {
	m := map[int]string{
		TypeLegislation: "law",
		TypeCase:        "case",
		TypeArticle:     "hotnews,ip_hottopic,ep_news_law,ep_news_case",
	}
	if found, ok := m[t]; ok {
		return found
	}
	return "law"
}

func getContentTypeText(ct int) string {
	m := map[int]string{
		TypeLegislation: "法规",
		TypeCase:        "案例",
		TypeArticle:     "评论文章",
	}
	if found, ok := m[ct]; ok {
		return found
	}
	return "法规"
}

func getTitles(keyword, db string) ([][]string, error) {
	autnResp := AutnResponse{}
	url := "http://192.168.2.210:9003/a=query&databasematch=%s" +
		"&sort=relevance+power_level:numberincreasing+date&print=none" +
		"&text=(%s):dretitle:articleid+OR+(%s):tags"
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
		article := []string{hit.Title, hit.Reference}
		articles = append(articles, article)
	}

	return articles, nil
}
