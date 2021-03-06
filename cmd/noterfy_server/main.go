package main

import (
	"log"
	"noterfy/api"
	"noterfy/api/middleware"
	"noterfy/api/server"
	"noterfy/api/server/routes"
	"noterfy/config"
	"noterfy/note/api/v1/transport/rest"
	noteservice "noterfy/note/service"
	filestore "noterfy/note/store/file"
	"os"
	"path/filepath"
	"time"
)

var (
	// Version is the version of the current server
	Version = "development"
	// BuildCommit is the git build recent commit during server build.
	BuildCommit = "development"
	// BuildDate is the timestamp of when the server last build.
	BuildDate = time.Now().Truncate(time.Second).UTC()
)

const (
	dbFileName = "note.pb"
)

func main() {

	conf := config.New()

	metadata := &routes.Metadata{
		Version:     Version,
		BuildCommit: BuildCommit,
		BuildDate:   BuildDate,
	}

	file, err := os.OpenFile(filepath.Join(conf.Store.File.Path, dbFileName), os.O_CREATE|os.O_RDWR, 0666)
	mustNoError(err)
	defer func() { _ = file.Close() }()

	svc := noteservice.New(filestore.New(file))
	srv := server.New(&server.Config{
		Port:     conf.Server.Port,
		Metadata: metadata,
		Middlewares: []api.NamedMiddleware{
			middleware.NewLoggingMiddleware(),
			middleware.NewRateLimitMiddleware(middleware.RateLimitConfig{
				DefaultExpirationTTL: time.Second,
				ExpireJobInterval:    time.Second,
				MaxBurst:             1,
			}),
			middleware.NewCORSMiddleware(nil),
		},
	})

	srv.AddRoutes(routes.Routes(metadata)...)
	srv.AddRoutes(rest.Routes(svc)...)
	mustNoError(srv.ListenAndServe())
}

func mustNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
