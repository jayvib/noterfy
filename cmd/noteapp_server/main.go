package main

import (
	"log"
	"noteapp/api"
	"noteapp/api/middleware"
	"noteapp/api/server"
	"noteapp/api/server/meta"
	"noteapp/config"
	"noteapp/note/api/v1/transport/rest"
	noteservice "noteapp/note/service"
	filestore "noteapp/note/store/file"
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

	file, err := os.OpenFile(filepath.Join(conf.Store.File.Path, dbFileName), os.O_CREATE|os.O_RDWR, 0666)
	mustNoError(err)
	defer func() { _ = file.Close() }()

	svc := noteservice.New(filestore.New(file))
	srv := server.New(&server.Config{
		Port: conf.Server.Port,
		Middlewares: []api.NamedMiddleware{
			middleware.NewLoggingMiddleware(),
		},
	})

	srv.AddRoutes(meta.Routes(&meta.Metadata{
		Version:     Version,
		BuildCommit: BuildCommit,
		BuildDate:   BuildDate,
	})...)
	srv.AddRoutes(rest.Routes(svc)...)

	defer srv.Close()
	mustNoError(srv.ListenAndServe())
}

func mustNoError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
