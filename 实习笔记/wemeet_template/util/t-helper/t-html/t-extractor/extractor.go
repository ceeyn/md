package extractor

// html图片文本等数据提取接口
type TExtractor interface {
	Parse(content string) (*THtmlElements, error)
}
