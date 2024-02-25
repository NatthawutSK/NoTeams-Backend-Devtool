package servers

import (
	"github.com/NatthawutSK/NoTeams-Backend/modules/files/filesHandler"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files/filesRepository"
	"github.com/NatthawutSK/NoTeams-Backend/modules/files/filesUsecase"
)

type IFilesModule interface {
	Init()
	Usecase() filesUsecase.IFilesUsecase
	Handler() filesHandler.IFileHandler
	Repository() filesRepository.IFilesRepository
}

type filesModule struct {
	*moduleFactory
	usecase filesUsecase.IFilesUsecase
	handler filesHandler.IFileHandler
}

func (m *moduleFactory) FilesModule() IFilesModule {
	repo := filesRepository.FilesRepository(m.s.db)
	usecase := filesUsecase.FilesUsecase(m.s.cfg, repo)
	handler := filesHandler.FileHandler(m.s.cfg, usecase)

	return &filesModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (f *filesModule) Init() {
	router := f.r.Group("/files")

	router.Get("/team/:team_id", f.mid.JwtAuth(), f.mid.AuthTeam(), f.handler.GetFilesTeam)
	// router.Post("/upload", f.handler.UploadFiles)
}

func (f *filesModule) Usecase() filesUsecase.IFilesUsecase          { return f.usecase }
func (f *filesModule) Handler() filesHandler.IFileHandler           { return f.handler }
func (f *filesModule) Repository() filesRepository.IFilesRepository { return f.usecase }
