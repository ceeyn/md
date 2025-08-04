package namebadge

import (
	"io/ioutil"
	"os"
)

//PathExists 路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	//isnotexist来判断，是不是不存在的错误
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

//GetDirFileName 获取文件目录名称
func GetDirFileName(dirPath string) []string {

	files, _ := ioutil.ReadDir(dirPath)

	fileNames := make([]string, 0, len(files))
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		fileNames = append(fileNames, f.Name())
	}

	return fileNames
}
