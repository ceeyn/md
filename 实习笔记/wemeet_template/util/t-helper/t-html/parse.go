package thtml

import (
	"bytes"
	"errors"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"strings"
)

type ParseFunc func(node *html.Node) error

// 默认解析函数
func tDefaultParser(node *html.Node) error {
	return nil
}

// html解析节点数据
type TParser struct {
	tree       *html.Node
	parseFuncs map[string]ParseFunc
}

// html数据解析函数
func NewTParser() *TParser {
	return &TParser{
		parseFuncs: make(map[string]ParseFunc),
	}
}

// Register
func (parser *TParser) Register(nodeTag string, parseFunc ParseFunc) {
	parser.parseFuncs[nodeTag] = parseFunc
}

// 解析html
func (parser *TParser) Parse(content string) error {
	if content == "" {
		return fmt.Errorf("content is empty")
	}

	pnode, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return err
	}
	parser.tree = pnode

	err = parser.parse(pnode)
	if err != nil {
		return err
	}
	return err
}

// 解析html
func (parser *TParser) parse(pnode *html.Node) error {
	if pnode == nil {
		return nil
	}

	var err error

	if pnode.Type == html.ErrorNode {
		err = errors.New("invalid html content")
	} else if pnode.Type == html.TextNode {
		err = parser.searchParseFunc(THtmlTagText)(pnode)
	} else if pnode.Type == html.ElementNode {
		if pnode.Data == THtmlTagImage {
			err = parser.searchParseFunc(THtmlTagImage)(pnode)
		} else if pnode.Data == THtmlTagVideo {
			err = parser.searchParseFunc(THtmlTagVideo)(pnode)
		} else if pnode.Data == THtmlTagHref {
			err = parser.searchParseFunc(THtmlTagHref)(pnode)
		} else if pnode.Data == THtmlTagSource {
			err = parser.searchParseFunc(THtmlTagSource)(pnode)
		}
	}
	if err != nil {
		return err
	}

	if pnode.FirstChild != nil {
		err = parser.parse(pnode.FirstChild)
	}
	if pnode.NextSibling != nil {
		err = parser.parse(pnode.NextSibling)
	}
	return err
}

// 获取解析函数
func (parser *TParser) searchParseFunc(nodeTag string) ParseFunc {
	if parseFunc, ok := parser.parseFuncs[nodeTag]; ok {
		return parseFunc
	}
	return tDefaultParser
}

// 构建html文本
func (parser *TParser) Build() (string, error) {
	var res bytes.Buffer
	writer := io.Writer(&res)
	err := html.Render(writer, parser.tree)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}
