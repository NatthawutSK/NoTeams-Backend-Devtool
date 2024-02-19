package main

import (
	"log"
	"os"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/databases"
	"github.com/NatthawutSK/NoTeams-Backend/servers"
)

func main() {
	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is required")
		}
		return os.Args[1]
	}())

	// s3 := s3Conn.S3Connect(cfg.S3())
	// _ = s3

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	servers.NewSever(cfg, db).Start()
}
