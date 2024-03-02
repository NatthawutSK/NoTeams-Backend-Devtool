package files

import "mime/multipart"

type FileReq struct {
	File           *multipart.FileHeader `form:"files"`
	Destination    string
	Extension      string
	OriginFilename string
	FileName       string
	ContentType    string
}

type DeleteFilesTeamReq struct {
	FileId string `json:"file_id" form:"file_id" validate:"required"`
}
