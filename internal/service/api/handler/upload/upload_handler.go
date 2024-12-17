package upload

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"im/internal/logger"
	"im/internal/response"
	"im/internal/service/api/config"
	"im/pkg/ossclient"
	"strings"
	"time"
)

type UploadHandler struct {
	ossClient ossclient.ImOssClient
}

func NewUploadHandler(router fiber.Router) *UploadHandler {
	cfg := config.Config.Oss
	handler := &UploadHandler{
		ossClient: ossclient.NewAliyunOssClient(
			cfg.Endpoint,
			cfg.AccessKeyId,
			cfg.AccessKeySecret,
			cfg.BucketName,
			cfg.Region,
			cfg.Host),
	}
	router.Post("/file/upload", handler.upload)
	return handler
}
func (h *UploadHandler) upload(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	// => *multipart.Form

	// Get all files from "documents" key:
	files := form.File["file"]
	// => []*multipart.FileHeader
	if len(files) == 0 {
		return c.JSON(response.FailWithMsg("upload file is empty"))
	}

	file := files[0]
	logger.Infof("--- upload file start, file:%s,size:%d ---",
		file.Filename, file.Size)

	fh, err := file.Open()
	if err != nil {
		logger.Errorf("upload file open error,err:%v", err)
		return c.JSON(response.FailWithMsg("upload file error"))
	}
	defer fh.Close()

	//get suffix
	suffix := ".jpg"
	if len(file.Filename) > 0 {
		index := strings.LastIndex(file.Filename, ".")
		if index > 0 {
			suffix = strings.ToLower(file.Filename[index:])
		}
	}

	now := time.Now()
	key := fmt.Sprintf("upload/im/%d/%d/%d/%d_%s",
		now.Year(), now.Month(), now.Day(), now.UnixMicro(), suffix)
	url, err := h.ossClient.UploadByReader(key, fh)
	if err != nil {
		logger.Errorf("upload file to oss error,err:%v", err)
		return c.JSON(response.FailWithMsg("upload file to oss error"))
	}
	return c.JSON(response.Success(url))
}
