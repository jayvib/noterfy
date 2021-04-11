package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"noterfy/api"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"
)

// New takes config for all the arguments tat the server needs and
// return a server instance.
func New(conf *Config) *Server {
	server := &Server{
		Port:        conf.Port,
		Middlewares: conf.Middlewares,
		HTTPRoutes:  conf.HTTPRoutes,
	}
	return server
}

// Config contains all the arguments that the server needs.
type Config struct {
	// Port is the port of the server
	Port int
	// Middlewares are the middlewares to be apply to the
	// handler.
	Middlewares []api.NamedMiddleware
	// HTTPRoutes are the API routes that will be register to the server.
	HTTPRoutes []api.Route
}

// Server is the wrapper for all the bootstrapping of a typical server.
type Server struct {
	Port        int
	Middlewares []api.NamedMiddleware
	server      *http.Server
	HTTPRoutes  []api.Route
	isInited    bool
}

func (s *Server) init() {
	router := mux.NewRouter()

	s.printInfo()

	for _, routes := range s.HTTPRoutes {
		router.Path(routes.Path()).Methods(routes.Method()).Handler(routes.Handler())
	}

	for _, mw := range s.Middlewares {
		router.Use(mw.Middleware)
	}

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.Port),
		Handler: router,
	}

	s.isInited = true
}

func (s *Server) printInfo() {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', tabwriter.TabIndent)
	defer func() { _ = w.Flush() }()

	writeToConsole := func(format string, a ...interface{}) {
		_, _ = fmt.Fprintf(w, format, a...)
	}

	writeToConsole("==============ROUTES================\n")
	for _, route := range s.HTTPRoutes {
		writeToConsole("👉️ %s\t%s\n", route.Method(), route.Path())
	}
	writeToConsole("\n")
	writeToConsole("============MIDDLEWARE==============\n")
	for _, mw := range s.Middlewares {
		writeToConsole("👉 %s\n", mw.Name)
	}
	writeToConsole("\n")
}

// ListenAndServe serves clients request by the server.
func (s *Server) ListenAndServe() (err error) {
	if !s.isInited {
		s.init()
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("\U0001F7E2 Server Started Listening on %s\n", s.server.Addr)
		serr := s.server.ListenAndServe()
		if serr != nil && serr != http.ErrServerClosed && err == nil {
			err = serr
		}
	}()

	<-done
	fmt.Println("🛑 Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Println("💯 Server Exited Properly")
	return
}

// Close closes the underlying server.
func (s *Server) Close() {
	if err := s.server.Close(); err != nil {
		logrus.Error(err)
	}
}

// AddRoutes takes routes to register in server.
func (s *Server) AddRoutes(routes ...api.Route) {
	s.HTTPRoutes = append(s.HTTPRoutes, routes...)
}
