package replacer

import (
	"errors"
	"golang.org/x/net/html"
	t_html "meeting_template/util/t-helper/t-html"
	"strings"
)

// html节点数据替换
type TSimpleReplacer struct {
}

var _ TReplacer = &TSimpleReplacer{}

// 批量替换html里的图片元素为统一图片
func (replacer *TSimpleReplacer) ReplaceImageSrcToSame(content, src string) (string, error) {
	var err error

	parser := t_html.NewTParser()
	parser.Register(t_html.THtmlTagImage, func(node *html.Node) error {
		attrs := node.Attr
		for idx, attr := range attrs {
			if attr.Key != THtmlAttrSrc {
				continue
			}
			attrs[idx].Val = src
		}
		return nil
	})

	err = parser.Parse(content)
	if err != nil {
		return "", err
	}
	return parser.Build()
}

// 批量替换html里的图片元素
func (replacer *TSimpleReplacer) ReplaceImageSrc(content string, targets []string) (string, error) {
	var err error

	parser := t_html.NewTParser()
	targetIdx := 0
	parser.Register(t_html.THtmlTagImage, func(node *html.Node) error {
		attrs := node.Attr
		for idx, attr := range attrs {
			// fmt.Println(idx, attr, targetIdx, len(targets))
			if attr.Key != THtmlAttrSrc {
				continue
			}
			// 新解析的图片元素比之前提取的多时，替换失败！
			if len(targets) <= targetIdx {
				return errors.New("parse fail or discription changed, new des has more image")
			}
			attrs[idx].Val = targets[targetIdx]
			targetIdx += 1
		}
		return nil
	})

	err = parser.Parse(content)
	if err != nil {
		return "", err
	}
	return parser.Build()
}

// 替换文本节点
func (replacer *TSimpleReplacer) ReplaceTextSrc(content string, targets []string) (string, error) {
	var err error

	parser := t_html.NewTParser()
	targetIdx := 0
	parser.Register(t_html.THtmlTagText, func(node *html.Node) error {
		textFrag := strings.TrimSpace(node.Data)
		if textFrag != "" {
			// 新解析的图片元素比之前提取的多时，替换失败！
			if len(targets) <= targetIdx {
				return errors.New("parse fail or discription changed, new des has more text")
			}
			node.Data = targets[targetIdx]
			targetIdx += 1
		}
		return nil
	})

	err = parser.Parse(content)
	if err != nil {
		return "", err
	}
	return parser.Build()
}
