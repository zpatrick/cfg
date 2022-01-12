package main

import (
	"context"
	"log"

	"github.com/zpatrick/cfg/example/config"
	"github.com/zpatrick/cfg/example/database"
	"github.com/zpatrick/cfg/example/server"
)

func main() {
	ctx := context.Background()
	conf, err := config.Load(ctx)
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.CreateDB(ctx, conf.DB)
	if err != nil {
		log.Fatal(err)
	}

	svr, err := server.CreateServer(ctx, db, conf.Server)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("running service on port:", conf.Server.Port)
	log.Fatal(svr.ListenAndServe())
}
