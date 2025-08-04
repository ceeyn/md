package extractor

import (
	"golang.org/x/net/html"
	t_html "meeting_template/util/t-helper/t-html"
	"strings"
)

// 数据提取实现接口
type TSimpleExtractor struct {
}

var _ TExtractor = &TSimpleExtractor{}

// html数据解析
func (extractor *TSimpleExtractor) Parse(content string) (*THtmlElements, error) {
	var err error
	var res = &THtmlElements{}

	parser := t_html.NewTParser()
	parser.Register(t_html.THtmlTagText, func(node *html.Node) error {
		textFrag := strings.TrimSpace(node.Data)
		if textFrag != "" {
			res.TextFrags = append(res.TextFrags, textFrag)
		}
		return nil
	})
	parser.Register(t_html.THtmlTagImage, func(node *html.Node) error {
		attrs := node.Attr
		for _, attr := range attrs {
			if attr.Key != THtmlAttrSrc {
				continue
			}
			res.Images = append(res.Images, attr.Val)
		}
		return nil
	})
	parser.Register(t_html.THtmlTagSource, func(node *html.Node) error {
		if node.Type == html.ElementNode && node.Data == t_html.THtmlTagSource {
			attrs := node.Attr
			for _, attr := range attrs {
				if attr.Key != THtmlAttrSrc {
					continue
				}
				res.Videos = append(res.Videos, attr.Val)
			}
		}
		return nil
	})

	err = parser.Parse(content)
	return res, err
}
