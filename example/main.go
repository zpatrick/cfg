package main

import (
	"context"
	"log"
	"os"

	"github.com/zpatrick/cfg/example/config"
	"github.com/zpatrick/cfg/example/db"
	"github.com/zpatrick/cfg/example/svr"
)

func main() {
	ctx := context.Background()
	conf, err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	server, err := svr.CreateServer(ctx, *conf.Server)
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.CreateDB(ctx, *conf.DB)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("service successfully configured")
	os.Exit(run(database, server))
}

func run(database *db.DB, server *svr.Server) int {
	return 0
}
