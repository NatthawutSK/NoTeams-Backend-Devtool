package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files"

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

type filesPub struct {
	bucket      string
	destination string
	file        *files.FileRes
}

func (f *filesPub) makePublic(ctx context.Context, client *storage.Client) error {
	acl := client.Bucket(f.bucket).Object(f.destination).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return fmt.Errorf("ACLHandle.Set: %w", err)
	}
	fmt.Printf("Blob %v is now publicly accessible.\n", f.destination)
	return nil
}

func (u *upload) uploadWorkers(ctx context.Context, client *storage.Client, jobs <-chan *files.FileReq, result chan<- *files.FileRes, errs chan<- error) {
	//jobs <-chan คือการรับค่าจาก channel แบบ receive only
	//result chan<- คือการส่งค่าไปที่ channel แบบ send only
	//errs chan<- คือการส่งค่าไปที่ channel แบบ send only

	for job := range jobs {
		container, err := job.File.Open()
		if err != nil {
			errs <- fmt.Errorf("open file failed: %v", err)
			return
		}
		b, err := io.ReadAll(container)
		if err != nil {
			errs <- fmt.Errorf("read file failed: %v", err)
			return
		}
		buf := bytes.NewBuffer(b)

		// Upload an object with storage.Writer.
		wc := client.Bucket(u.cfg.App().GCPBucket()).Object(job.Destination).NewWriter(ctx)
		wc.ObjectAttrs.ContentType = job.ContentType

		if _, err = io.Copy(wc, buf); err != nil {
			errs <- fmt.Errorf("io.Copy: %w", err)
			return
		}
		// Data can continue to be added to the file until the writer is closed.
		if err := wc.Close(); err != nil {
			errs <- fmt.Errorf("Writer.Close: %w", err)
			return
		}
		fmt.Printf("%v uploaded to %v.\n", job.FileName, job.Destination)

		newFile := &filesPub{
			file: &files.FileRes{
				FileName: job.OriginFilename,
				Url:      fmt.Sprintf("https://storage.googleapis.com/%s/%s", u.cfg.App().GCPBucket(), job.Destination),
			},
			bucket:      u.cfg.App().GCPBucket(),
			destination: job.Destination,
		}

		if err := newFile.makePublic(ctx, client); err != nil {
			errs <- fmt.Errorf("make file public failed: %v", err)
			return
		}

		errs <- nil
		result <- newFile.file
	}

}

func (u *upload) UploadFiles(filesReq []*multipart.FileHeader, isDownload bool, folder string) ([]*files.FileRes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	fmt.Println("uploading files to GCP bucket : ", u.cfg.App().GCPBucket())

	filesUpload := make([]*files.FileReq, 0)
	res := make([]*files.FileRes, 0)
	contentType := "application/octet-stream"

	// files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
		"pdf":  "pdf",
	}

	for _, file := range filesReq {
		if !isDownload {
			contentType = file.Header.Get("Content-Type")
		}
		// check file extension
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return nil, fmt.Errorf("invalid filesReq extension")
		}

		// check filesReq size
		if file.Size > int64(u.cfg.App().FileLimit()) {
			return nil, fmt.Errorf("filesReq size must less than %d MB", int(math.Ceil(float64(u.cfg.App().FileLimit())/math.Pow(1024, 2))))
		}

		filename := RandFileName(ext)
		// if folder != "" {
		// 	filename = fmt.Sprintf("%s/%s", folder, filename)
		// }
		fileUp := &files.FileReq{
			File:           file,
			FileName:       filename,
			OriginFilename: file.Filename,
			Destination:    fmt.Sprintf("%s/%s", folder, filename),
			Extension:      ext,
			ContentType:    contentType,
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
		go u.uploadWorkers(ctx, client, jobsCh, resultsCh, errorsCh)
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
