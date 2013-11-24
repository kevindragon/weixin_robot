// 解析输入字符串，把指令和关键词提取出来
package main

import (
	"strings"
)

const (
	CmdNone = 0
	CmdHelp = 1 << (iota - 1)
	CmdAccountBind
	CmdSearch
)

// 解析从微信发过来的字符串
func parseCmd(input []byte) int {
	inputStr := string(input)

	// 帮助命令
	if inputStr == "帮助" || inputStr == "?" || inputStr == "？" {
		return CmdHelp
	}
	if inputStr == "账号绑定" || inputStr == "绑定账号" {
		return CmdAccountBind
	}
	// 默认为搜索
	return CmdSearch
}

func parseSearchCmd(input []byte) (int, string) {
	inputStr := string(input)
	inputRune := []rune(inputStr)
	ct := TypeNone
	if len(inputRune) >= 3 {
		ct = getContentType(string(inputRune[:3]))
	}
	keyword := strings.Trim(string(inputRune), " ")
	if ct != TypeNone && len(inputRune) >= 3 {
		keyword = strings.Trim(string(inputRune[3:]), " ")
	}
	return ct, keyword
}
