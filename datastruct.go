package main

import (
	"encoding/xml"
)

// 文本消息基本内容
type TextMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	Content      string
}

// 接收到的消息
type TextMessageReceived struct {
	TextMessage
	MsgId string
}

// 图文消息
type NewsMessage struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string
	FromUserName string
	CreateTime   int64
	MsgType      string
	ArticleCount int
	Articles     NewsMessageArtice
}
type NewsMessageArtice struct {
	XMLName xml.Name `xml:"Articles"`
	Item    []NewsMessageItem
}
type NewsMessageItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string
	Description string
	PicUrl      string
	Url         string
}

type AutnResponse struct {
	Action   string       `xml:"action"`
	Response string       `xml:"response"`
	Respdata Responsedata `xml:"responsedata"`
}
type Responsedata struct {
	Numhits int   `xml:"numhits"`
	Hits    []Hit `xml:"hit"`
}
type Hit struct {
	Reference string `xml:"reference"`
	Title     string `xml:"title"`
}
