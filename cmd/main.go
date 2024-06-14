package main

import (
	"context"
	"entrypoint/config"
	"entrypoint/entrypoint"
	"entrypoint/storage/mongodb"
	"github.com/rs/zerolog/log"
	"time"
)

func main() {
	cfg, _ := config.NewConfig("etc/config.yaml")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	mongo, err := mongodb.NewClient(cfg.MongoDB.Url)

	if err != nil {
		log.Fatal().Err(err).Str("func", "main").Msg("Failed create mongodb client")
	}

	defer cancel()
	defer mongo.Disconnect(ctx)

	defaultHandler := &entrypoint.DefaultHandler{
		MongoClient: mongo,
		Db:          cfg.MongoDB.Db,
		Coll:        cfg.MongoDB.Coll,
	}

	newServer := entrypoint.NewServer(cfg.Entrypoint.Host)
	newServer.RegisterHandler("/", &entrypoint.HandlerRoot{})
	newServer.RegisterHandler("/api/v1", &entrypoint.HandlerV1Api{})
	newServer.RegisterHandler("/api/v1/backends", defaultHandler)
	newServer.ServerStart()
}
