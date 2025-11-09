package main

import (
	"log"

	"search-service/internal/config"
	"search-service/internal/graphql"
	"search-service/internal/group"
	"search-service/internal/server"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	groupClient := group.NewHTTPClient(conf.GroupServiceURL, conf.HTTPClientTimeout)
	groupService := group.NewService(groupClient, conf.DefaultResultSize, conf.MaxResultSize)

	schemaBuilder := graphql.SchemaBuilder{GroupService: groupService}
	schema, err := schemaBuilder.Build()
	if err != nil {
		log.Fatalf("failed to build GraphQL schema: %v", err)
	}

	apiServer := server.New(conf, schema)
	if err := apiServer.Run(); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
