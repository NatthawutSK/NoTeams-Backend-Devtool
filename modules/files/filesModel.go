package files

import "mime/multipart"

type FileReq struct {
	Files       *multipart.FileHeader `form:"files"`
	FileName    string                `json:"file_name" form:"file_name"`
	ContentType string
}

type DeleteFilesTeamReq struct {
	FileId string `json:"file_id" form:"file_id" validate:"required"`
}
