package filesRepository

import (
	"encoding/json"
	"fmt"

	"github.com/NatthawutSK/NoTeams-Backend/modules/files"
	"github.com/jmoiron/sqlx"
)

type IFilesRepository interface {
	GetFilesTeam(teamId string) (*files.GetFilesTeamRes, error)
	// UploadFiles() error
}

type filesRepository struct {
	db *sqlx.DB
}

func FilesRepository(db *sqlx.DB) IFilesRepository {
	return &filesRepository{
		db: db,
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

	// var filesTeam []*files.GetFilesTeamRes
	// if err := r.db.Select(&filesTeam, query, teamId); err != nil {
	// 	return nil, fmt.Errorf("error: %w", err)
	// }
	// return filesTeam, nil
}

// func (r *filesRepository) UploadFiles() error {

// 	query := `
// 	INSERT INTO File (file_name, file_url, team_id, user_id)
// 	VALUES ($1, $2, $3, $4)
// 	`

// 	return nil
// }
