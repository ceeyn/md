package rpc

import (
	"io"
	"os"
	"path"
	"strings"

	"archive/zip"
	"git.code.oa.com/trpc-go/trpc-go/log"
)

//Decompress zip解压
func Decompress(zipFile, dest string) error {
	reader, err := zip.OpenReader(zipFile)
	defer reader.Close()
	if err != nil {
		return err
	}

	if dest != "" {
		if err = os.MkdirAll(dest, 0755); err != nil {
			return err
		}
	}

	for _, file := range reader.File {
		log.Infof("Decompress, zipfile:%+v, dest%+v, filename:%+v", zipFile, dest, file.Name)

		if filterInvalidFile(file.Name) {
			log.Infof("filterInvalidFile, filename:%+v", file.Name)
			continue
		}

		filename := path.Join(dest, file.Name)
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(filename, 0755)
			if err != nil {
				return err
			}
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}

		w, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			rc.Close()
			continue
		}

		_, err = io.Copy(w, rc)
		if err != nil {
			rc.Close()
			w.Close()
			return err
		}
		w.Close()
		rc.Close()
	}
	return nil
}

func filterInvalidFile(fileName string) bool {
	bFileName := []byte(fileName)
	if len(bFileName) == 0 {
		return true
	}
	if string(bFileName[0]) == "." {
		return true
	}
	if strings.Contains(fileName, "/.") {
		return true
	}
	return false
}
