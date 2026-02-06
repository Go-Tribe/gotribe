// Copyright 2024 Innkeeper GoTribe <https://www.gotribe.cn>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://www.gotribe.cn

package upload

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"strings"
)

// Provider 上传服务提供商
type Provider string

const (
	ProviderOSS   Provider = "oss"
	ProviderQiniu Provider = "qiniu"
	ProviderS3    Provider = "s3" // 预留：亚马逊 S3
)

// Options 上传服务配置，根据 Provider 使用不同字段（如 S3 需 Region）
type Options struct {
	Provider       Provider
	Endpoint       string
	AccessKey      string
	SecretKey      string
	Bucket         string
	Region         string // S3 等需要，可选
}

// Uploader 定义上传接口
type Uploader interface {
	UploadFile(file *multipart.FileHeader) (UploadResource, error)
	DeleteFile(key string) error
}

type UploadResource struct {
	FileExt string `json:"file_ext"`
	Key     string `json:"key"`
	Domain  string `json:"domain"`
}

// Service 提供上传和删除文件的功能
type Service struct {
	uploader Uploader
}

// NewService 根据 Options 创建对应的上传 Service（oss / qiniu，后续可扩展 s3）
func NewService(opts *Options) (*Service, error) {
	if opts == nil {
		return nil, errors.New("upload options is nil")
	}
	provider := Provider(strings.ToLower(string(opts.Provider)))
	var uploader Uploader
	switch provider {
	case ProviderOSS:
		uploader = NewOSS(opts.Endpoint, opts.AccessKey, opts.SecretKey, opts.Bucket)
	case ProviderQiniu:
		uploader = NewQiniu(opts.AccessKey, opts.SecretKey, opts.Bucket)
	case ProviderS3:
		return nil, fmt.Errorf("upload provider %q not implemented yet", provider)
	default:
		return nil, fmt.Errorf("unknown upload provider: %q", opts.Provider)
	}
	return &Service{uploader: uploader}, nil
}

// NewUploadFile 根据 endpoint/accessKey/secretKey/bucket 与 enableOss 创建 Service（兼容旧配置，建议改用 NewService + upload-file.provider）
func NewUploadFile(endpoint, accessKeyId, accessKeySecret, bucketName string, enableOss bool) (*Service, error) {
	p := ProviderQiniu
	if enableOss {
		p = ProviderOSS
	}
	return NewService(&Options{
		Provider:  p,
		Endpoint:  endpoint,
		AccessKey: accessKeyId,
		SecretKey: accessKeySecret,
		Bucket:    bucketName,
	})
}

// UploadFile 公用上传文件方法
func (s *Service) UploadFile(file *multipart.FileHeader) (UploadResource, error) {
	if s.uploader == nil {
		return UploadResource{}, os.ErrInvalid
	}
	result, err := s.uploader.UploadFile(file)
	if err != nil {
		return UploadResource{}, err
	}
	return result, nil
}

// DeleteFile 公用删除文件方法
func (s *Service) DeleteFile(key string) error {
	if s.uploader == nil {
		return os.ErrInvalid
	}
	err := s.uploader.DeleteFile(key)
	if err != nil {
		return err
	}
	return nil
}
