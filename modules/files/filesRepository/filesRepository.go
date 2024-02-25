package filesRepository

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/NatthawutSK/NoTeams-Backend/modules/files"
	"github.com/jmoiron/sqlx"
)

type IFilesRepository interface {
	GetFilesTeam(teamId string) (*files.GetFilesTeamRes, error)
	UploadFilesTeam(userId string, teamId string, req []*files.FileRes) ([]*files.FileTeamByIdRes, error)
}

type filesRepository struct {
	db   *sqlx.DB
	pCtx context.Context
}

func FilesRepository(db *sqlx.DB, pCtx context.Context) IFilesRepository {
	return &filesRepository{
		db:   db,
		pCtx: pCtx,
	}
}

func (r *filesRepository) GetFilesTeam(teamId string) (*files.GetFilesTeamRes, error) {

	query := `
	SELECT
		COALESCE(array_to_json(array_agg("files")), '[]'::json)
	FROM (
		SELECT
			"f"."file_name",
			"f"."file_url",
			"f"."created_at",
			"u"."username"
		from "File" f
		JOIN "User" u ON f."user_id" = u."user_id"
		WHERE f."team_id" = $1
	) AS "files";
	`

	FilesBytes := make([]byte, 0)
	if err := r.db.Get(&FilesBytes, query, teamId); err != nil {
		return nil, fmt.Errorf("get files team failed: %v", err)
	}

	files := make(files.GetFilesTeamRes, 0)
	if err := json.Unmarshal(FilesBytes, &files); err != nil {
		return nil, fmt.Errorf("failed to unmarshal files team: %v", err)
	}

	return &files, nil

}

func (r *filesRepository) UploadFilesTeam(userId string, teamId string, req []*files.FileRes) ([]*files.FileTeamByIdRes, error) {
	ctx, cancel := context.WithTimeout(r.pCtx, 20*time.Second)
	defer cancel()

	filesId := make([]string, 0)

	queryFiles := `
		INSERT INTO "File" (
			team_id,
			user_id,
			file_name,
			file_url
			)
		VALUES ($1, $2, $3, $4)
		RETURNING "file_id";
		`
	for _, file := range req {
		var id string
		filename := filepath.Base(file.FileName)
		if err := r.db.QueryRowContext(
			ctx,
			queryFiles,
			teamId,
			userId,
			filename,
			file.Url,
		).Scan(&id); err != nil {
			return nil, fmt.Errorf("insert oauth failed: %v", err)
		}

		filesId = append(filesId, id)
	}

	query := `
		SELECT
			"f"."file_name",
			"f"."file_url",
			"f"."created_at",
			"u"."username"
		FROM "File" f
		JOIN "User" u ON f."user_id" = u."user_id"
		WHERE f."file_id" = $1`

	filesRes := make([]*files.FileTeamByIdRes, 0)

	for _, fileId := range filesId {

		file := new(files.FileTeamByIdRes)

		if err := r.db.Get(file, query, fileId); err != nil {
			return nil, fmt.Errorf("get files by file id failed: %v", err)
		}

		fmt.Println("in repo", file)

		filesRes = append(filesRes, file)

	}

	return filesRes, nil
}
