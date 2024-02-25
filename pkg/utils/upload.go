package utils

import (
	"context"
	"fmt"
	"math"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/modules/files"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/s3Conn"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/NatthawutSK/NoTeams-Backend/config"
)

type IUpload interface {
	UploadFiles(req []*multipart.FileHeader, isDownload bool, folder string) ([]*files.FileRes, error)
}

type upload struct {
	cfg config.IConfig
}

func Upload(cfg config.IConfig) IUpload {
	return &upload{
		cfg: cfg,
	}
}

func (u *upload) UploadFiles(filesReq []*multipart.FileHeader, isDownload bool, folder string) ([]*files.FileRes, error) {
	s3Client := s3Conn.S3Connect(u.cfg.S3())
	contentType := "application/octet-stream"
	filesUpload := make([]*files.FileReq, 0)
	res := make([]*files.FileRes, 0)

	// files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
		"pdf":  "pdf",
	}

	for _, file := range filesReq {
		// check file extension
		if !isDownload {
			contentType = file.Header.Get("Content-Type")
		}
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return nil, fmt.Errorf("invalid filesReq extension")
		}

		// check filesReq size
		if file.Size > int64(u.cfg.App().FileLimit()) {
			return nil, fmt.Errorf("filesReq size must less than %d MB", int(math.Ceil(float64(u.cfg.App().FileLimit())/math.Pow(1024, 2))))
		}

		filename := RandFileName(ext)
		if folder != "" {
			filename = fmt.Sprintf("%s/%s", folder, filename)
		}
		fileUp := &files.FileReq{
			FileName:    filename,
			Files:       file,
			ContentType: contentType,
		}

		filesUpload = append(filesUpload, fileUp)
	}

	jobsCh := make(chan *files.FileReq, len(filesUpload))
	resultsCh := make(chan *files.FileRes, len(filesUpload))
	errorsCh := make(chan error, len(filesUpload))

	for _, r := range filesUpload {
		jobsCh <- r
	}
	close(jobsCh)

	numWorkers := 5
	for i := 0; i < numWorkers; i++ {
		go u.uploadWorkers(s3Client, jobsCh, resultsCh, errorsCh)
	}

	for a := 0; a < len(filesUpload); a++ {
		err := <-errorsCh
		if err != nil {
			return nil, fmt.Errorf("upload file failed: %v", err)
		}
		result := <-resultsCh
		res = append(res, result)
	}

	return res, nil

}

func (u *upload) uploadWorkers(s3Client *s3.Client, jobs <-chan *files.FileReq, result chan<- *files.FileRes, errs chan<- error) {

	for job := range jobs {
		f, err := job.Files.Open()
		if err != nil {
			errs <- fmt.Errorf("open file failed: %v", err)
			return
		}
		defer f.Close()

		input := &s3.PutObjectInput{
			Bucket:      aws.String(u.cfg.S3().S3Bucket()),
			Key:         aws.String(job.FileName),
			Body:        f,
			ContentType: aws.String(job.ContentType),
		}

		_, err = s3Client.PutObject(context.TODO(), input)
		if err != nil {
			errs <- fmt.Errorf("put object failed: %v", err)
			return
		}

		newFile := &files.FileRes{
			FileName: job.FileName,
			Url:      fmt.Sprintf("https://%s.s3.amazonaws.com/%s", u.cfg.S3().S3Bucket(), job.FileName),
		}

		errs <- nil
		result <- newFile
	}

}
