package hdfs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	hdfs2 "github.com/colinmarc/hdfs/v2"
)

type (
	Client struct {
		Address    string
		User       string
		hdfsClient *hdfs2.Client
	}

	File struct {
		Name  string
		IsDir bool
		Files []File
	}
)

// NewHdfsClient 初始化客户端
func NewHdfsClient(address string, user string) (*Client, error) {
	hdfsClient, err := hdfs2.NewClient(hdfs2.ClientOptions{Addresses: strings.Split(address, ","), User: user})
	if err != nil {
		return nil, err
	}

	return &Client{
		Address:    address,
		User:       user,
		hdfsClient: hdfsClient,
	}, nil
}

// StatPath 获取文件或者路径的状态
func (hs *Client) StatPath(path string) (os.FileInfo, error) {
	return hs.hdfsClient.Stat(path)
}

// DownloadDirectory 下载某一个目录
func (hs *Client) DownloadDirectory(src string, dest string) error {
	fileName := filepath.Base(src)
	return hs.hdfsClient.Walk(src,
		func(path string, f os.FileInfo, err error) error {
			localPath := dest + "/"
			value := strings.Split(path, "/")
			for k, v := range value {
				// 此处的代码主要用于目标目录的文件，例如： /aw_data/simulation_offline/scene/test5_1665312743/dynamic_obstacle/32768_346173/15988
				if v == fileName {
					for key, val := range value[k+1:] {
						if key == len(value[k+1:])-1 {
							localPath += val
						} else {
							localPath += val + "/"
						}
					}

					break
				}
			}

			if f.IsDir() {
				_, err := os.Stat(localPath)
				if err != nil {
					if os.IsNotExist(err) {
						return os.Mkdir(localPath, os.ModePerm)
					}
					return err
				}
				return os.Chmod(localPath, os.ModePerm)
			}

			err = hs.DownloadFile(path, localPath)
			if err != nil {
				return err
			}

			return nil
		})
}

// DownloadFile 下载文件
func (hs *Client) DownloadFile(src, destPath string) (err error) {
	return hs.hdfsClient.CopyToLocal(src, destPath)
}

// DeleteFile 删除文件
func (hs *Client) DeleteFile(path string) (err error) {
	return hs.hdfsClient.RemoveAll(path)
}

// CreateDir 创建文件夹
func (hs *Client) CreateDir(uploadPath string) (err error) {
	return hs.hdfsClient.Mkdir(uploadPath, os.ModePerm)
}

// UploadFile 上传文件
func (hs *Client) UploadFile(src, sourcePath string) (err error) {
	err = hs.hdfsClient.CopyToRemote(src, sourcePath)
	if err != nil {
		return err
	}

	return hs.hdfsClient.Chmod(sourcePath, os.ModePerm)
}

// Rename 重命名
func (hs *Client) Rename(src, source string) (err error) {
	return hs.hdfsClient.Rename(src, source)
}

// ForeachDir 遍历文件夹
func (hs *Client) ForeachDir(src string) ([]os.FileInfo, error) {
	files, err := hs.hdfsClient.ReadDir(src)
	if err != nil {
		return nil, err
	}

	var myFile []os.FileInfo
	for _, file := range files {
		if file.IsDir() {
			myFile = append(myFile, file)
			continue
		}
	}

	return myFile, err
}

// ParseJson 解析数据
func (hs *Client) ParseJson(filePath string, v interface{}) error {
	data, err := hs.hdfsClient.ReadFile(filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &v)
}

// Files 获取目录下的n级文件
func (hs *Client) Files(dir string, deep uint) ([]File, error) {
	if deep == 0 {
		return nil, nil
	}

	list, err := hs.hdfsClient.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	deep--
	res := make([]File, 0, len(list))

	for _, each := range list {
		if each.IsDir() {
			if deep > 0 {
				f := File{Name: each.Name(), IsDir: true}
				f.IsDir = true
				f.Files, err = hs.Files(filepath.Join(dir, each.Name()), deep)
				if err != nil {
					return nil, err
				}
				res = append(res, f)
			}
		} else {
			res = append(res, File{Name: each.Name()})
		}
	}

	return res, nil
}
