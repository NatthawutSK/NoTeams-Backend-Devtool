package main

import (
	"fmt"
	"os"

	"github.com/NatthawutSK/NoTeams-Backend/config"
	"github.com/NatthawutSK/NoTeams-Backend/pkg/databases"
	"github.com/NatthawutSK/NoTeams-Backend/servers"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	cfg := config.LoadConfig(envPath())

	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	fmt.Println(db)

	servers.NewSever(cfg, db).Start()
}
