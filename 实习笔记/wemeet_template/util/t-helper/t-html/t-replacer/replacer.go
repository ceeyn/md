package replacer

// html数据节点替换接口
type TReplacer interface {
	ReplaceImageSrcToSame(content, src string) (string, error)
	ReplaceImageSrc(content string, targets []string) (string, error)
	ReplaceTextSrc(content string, targets []string) (string, error)
}
