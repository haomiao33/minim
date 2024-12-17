package ossclient

import "io"

type ImOssClient interface {
	Upload(key string, pathFile string) (string, error)
	UploadByReader(key string, reader io.Reader) (string, error)
}
