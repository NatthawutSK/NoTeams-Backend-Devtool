package files

import "mime/multipart"

type FileReq struct {
	Files       *multipart.FileHeader `form:"files"`
	FileName    string                `json:"file_name" form:"file_name"`
	ContentType string
}

type UploadFilesTeamReq struct{}
