package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"time"

	"git.code.oa.com/trpc-go/trpc-go/log"
	"git.code.oa.com/trpc-go/trpc-go/metrics"
	pb "git.code.oa.com/trpcprotocol/wemeet/common_upload"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
)

//URLToken ..
type URLToken struct {
	SessionToken string `url:"x-cos-security-token,omitempty" header:"-"`
}

// QueryCosTempKey 获取cos临时密钥信息
func QueryCosTempKey(ctx context.Context) (*pb.TempKeyInfo, error) {

	req := &pb.TempKeyV2Req{
		BusinessId: "wemeet_webinar",
	}

	proxy := pb.NewCommonUploadTrpcClientProxy()
	rsp, err := proxy.TrpcTempKeyV2(ctx, req)
	if err != nil {
		metrics.IncrCounter("TrpcTempKeyV2 failed", 1)
		log.ErrorContext(ctx, "[QueryCosTempKey] TrpcTempKeyV2 failed, req:%+v, err:%v", req, err)
		return nil, err
	}

	metrics.IncrCounter("TrpcTempKeyV2 success", 1)
	log.InfoContextf(ctx, "[QueryCosTempKey] TrpcTempKeyV2 success, req:%+v, rsp:%+v", req, rsp)
	return rsp.GetMajor(), nil
}

// DownLoadingFileForHttpUrl 通过url下载文件
func DownLoadingFileForHttpUrl(ctx context.Context, srcName, srcPath, fileUrl string) (string, error) {

	// 此处需要校验url
	if len(fileUrl) == 0 {
		metrics.IncrCounter("error fileUrl", 1)
		log.ErrorContextf(ctx, "[DownLoadingFileForHttpUrl] error fileUrl")
		return "", errors.New("error fileUrl")
	}

	rsp, err := http.Get(fileUrl)
	if err != nil {
		metrics.IncrCounter("http get failed", 1)
		log.ErrorContextf(ctx, "[DownLoadingFileForHttpUrl] http get failed, err:%v", err)
		return "", err
	}
	defer rsp.Body.Close()

	err = os.MkdirAll(srcPath, 755)
	if err != nil {
		metrics.IncrCounter("create MkdirAll failed", 1)
		log.ErrorContextf(ctx, "[DownLoadingFileForHttpUrl] create MkdirAll failed, err:%v, fileName", err, srcPath)
		return "", err
	}

	filePath := path.Join(srcPath, srcName) + ".zip"
	f, err := os.Create(filePath)
	if err != nil {
		metrics.IncrCounter("create file failed", 1)
		log.ErrorContextf(ctx, "[DownLoadingFileForHttpUrl] create file failed, err:%v, fileName", err, filePath)
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, rsp.Body)
	if err != nil {
		metrics.IncrCounter("copy failed", 1)
		log.ErrorContextf(ctx, "[DownLoadingFileForHttpUrl] copy failed, err:%v", err)
		return filePath, err
	}

	metrics.IncrCounter("DownLoadingFileForHttpUrl Success", 1)
	log.InfoContextf(ctx, "[DownLoadingFileForHttpUrl] success. srcName:%v, srcPath:%v, "+
		"fileUrl:%v, fileName:%v, ", srcName, srcPath, fileUrl, filePath)
	return filePath, nil
}

//UploadingFile ..
func UploadingFile(ctx context.Context, srcName, srcPath string) (string, error) {
	//获取临时秘钥
	tmpInfo, err := QueryCosTempKey(ctx)
	if err != nil {
		metrics.IncrCounter("QueryCosTempKey fail", 1)
		log.ErrorContextf(ctx, "QueryCosTempKey fail, err:%+v", err)
		return "", err
	}

	log.InfoContextf(ctx, "UploadingFile, name:%v,path:%v", srcName, srcPath)

	//文件上传
	u, _ := url.Parse(tmpInfo.GetBucketUrl())
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:     tmpInfo.GetTempSecretId(),
			SecretKey:    tmpInfo.GetTempSecretKey(),
			SessionToken: tmpInfo.GetSessionToken(),
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	key := fmt.Sprintf("%s", srcName)
	_, _, err = c.Object.Upload(ctx, key, path.Join(srcPath, srcName), nil)
	if err != nil {
		metrics.IncrCounter("Pic Upload fail", 1)
		log.ErrorContextf(ctx, "cos upload err:%+v", err.Error())
		return "", err
	}
	name := fmt.Sprintf("%s", srcName)
	token := &URLToken{
		SessionToken: tmpInfo.GetSessionToken(),
	}

	preSignedURL, err := c.Object.GetPresignedURL(ctx, http.MethodGet, name, tmpInfo.GetTempSecretId(),
		tmpInfo.GetTempSecretKey(), time.Minute*60, token)
	if err != nil {
		metrics.IncrCounter("cos presigned fail", 1)
		log.ErrorContextf(ctx, "cos presigned err:%+v", err.Error())
		return "", err
	}
	metrics.IncrCounter("UploadingFile succeed", 1)
	log.InfoContextf(ctx, "UploadingFile succeed and url is [%v]", preSignedURL.String())
	return preSignedURL.String(), nil
}
