package filesUsecase

import (
	"context"
	"fmt"
	"math"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/s3Conn"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type IFilesUsecase interface {
	UploadFiles(req []*multipart.FileHeader) ([]*files.FileRes, error)
	// UploadFile(client *s3.Client, bucket, filename string, fileHeader *multipart.FileHeader) (string, error)
}

type filesUsecase struct {
	cfg config.IConfig
}

func FilesUsecase(cfg config.IConfig) IFilesUsecase {
	return &filesUsecase{
		cfg: cfg,
	}
}

func (u *filesUsecase) UploadFiles(filesReq []*multipart.FileHeader) ([]*files.FileRes, error) {
	s3Client := s3Conn.S3Connect(u.cfg.S3())

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
		contentType := file.Header.Get("Content-Type")
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return nil, fmt.Errorf("invalid filesReq extension")
		}

		// check filesReq size
		if file.Size > int64(u.cfg.App().FileLimit()) {
			return nil, fmt.Errorf("filesReq size must less than %d MB", int(math.Ceil(float64(u.cfg.App().FileLimit())/math.Pow(1024, 2))))
		}

		filename := utils.RandFileName(ext)
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

// func (u *filesUsecase) UploadFile(client *s3.Client, bucket, filename string, fileHeader *multipart.FileHeader) (string, error) {
// 	// Open the file associated with the file header
// 	file, err := fileHeader.Open()
// 	if err != nil {
// 		return "", err
// 	}
// 	defer file.Close()

// 	// Upload the file to S3
// 	input := &s3.PutObjectInput{
// 		Bucket:      aws.String(bucket),
// 		Key:         aws.String(filename),
// 		Body:        file,
// 		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
// 	}

// 	_, err = client.PutObject(context.TODO(), input)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Construct the URL of the uploaded object
// 	url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucket, filename)
// 	return url, nil
// }

func (u *filesUsecase) uploadWorkers(s3Client *s3.Client, jobs <-chan *files.FileReq, result chan<- *files.FileRes, errs chan<- error) {

	for job := range jobs {
		f, err := job.Files.Open()
		if err != nil {
			errs <- fmt.Errorf("open file failed: %v", err)
			return
		}
		defer f.Close()

		// largeObject, err := io.ReadAll(f)
		// if err != nil {
		// 	errs <- fmt.Errorf("read file failed: %v", err)
		// 	return
		// }

		// largeBuffer := bytes.NewReader(largeObject)

		// var partMiBs int64 = 10
		// uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		// 	u.PartSize = partMiBs * 1024 * 1024
		// })
		// uploader := manager.NewUploader(s3Client)
		// _, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		// 	Bucket: aws.String(bucketName),
		// 	Key:    aws.String(objectKey),
		// 	Body:   largeBuffer,
		// })
		// if err != nil {
		// 	errs <- fmt.Errorf("put object failed: %v", err)
		// 	return
		// }

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
		// Construct the URL of the uploaded object
		// url := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", u.cfg.S3().S3Bucket(), filename)
	}

}