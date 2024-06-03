// Copyright 2023 Innkeeper GoTribe <https://ww.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a Apache style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/dengmengmian/solocms

package qiniu

import (
	"context"
	"mime/multipart"
	"path"
	"strconv"
	"time"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type FileRet struct {
	URL     string `json:"url"`
	FileExt string `json:"fileExt"`
	Key     string `json:"key"`
}

func UploadFile(file *multipart.FileHeader, ak, sk, bucket, domain string) (*FileRet, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	putPlicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(ak, sk)

	// 获取上传凭证
	upToken := putPlicy.UploadToken(mac)

	// 配置参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 华南区
		UseCdnDomains: false,
		UseHTTPS:      false, // 非https
	}
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数

	// 上传 自定义key，可以指定上传目录及文件名和后缀，
	currentTime := time.Now().Format("20060102")
	fileUnixName := strconv.FormatInt(time.Now().UnixNano(), 10)
	fileExt := path.Ext(file.Filename)
	key := currentTime + "/" + fileUnixName + fileExt // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	err = formUploader.Put(context.Background(), &ret, upToken, key, src, file.Size, &putExtra)

	if err != nil {
		return nil, err
	}

	url := domain + ret.Key // 返回上传
	return &FileRet{
		URL:     url,
		FileExt: fileExt,
		Key:     key,
	}, err
}

// 删除文件
// https://developer.qiniu.com/kodo/1238/go
func DeletdFile(accessKey, secretKey, bucket, key string) error {
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(bucket, key)
	if err != nil {
		return err
	}
	return nil
}
