package server

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"noterfy/api"
	"noterfy/api/server/meta"
	"os"
	"os/signal"
	"syscall"
	"text/tabwriter"
	"time"
)

const defaultShutdownTimeout = 5 * time.Second

// New takes config for all the arguments tat the server needs and
// return a server instance.
func New(conf *Config) *Server {
	conf.checkDefaults()
	server := &Server{
		Port:            conf.Port,
		Middlewares:     conf.Middlewares,
		HTTPRoutes:      conf.HTTPRoutes,
		ShutdownTimeout: conf.ShutdownTimeout,
		Metadata:        conf.Metadata,
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
	// ShutdownTimeout is at the duration to wait before abandoning the server's
	// shutdown. Default is 5 seconds.
	ShutdownTimeout time.Duration
	// Metadata is the API extra information.
	Metadata *meta.Metadata
}

func (c *Config) checkDefaults() {
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = defaultShutdownTimeout
	}
}

// Server is the wrapper for all the bootstrapping of a typical server.
type Server struct {
	Port        int
	Middlewares []api.NamedMiddleware
	server      *http.Server
	HTTPRoutes  []api.Route
	isInited    bool
	// ShutdownTimeout is at the duration to wait before abandoning the server's
	// shutdown. Default is 5 seconds.
	ShutdownTimeout time.Duration
	Metadata        *meta.Metadata
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

	writeToConsole("===============INFO=================\n")
	for _, v := range [][2]interface{}{
		{"Shutdown Timeout:\t%s\n", s.ShutdownTimeout},
		{"Version:\t%v\n", s.Metadata.Version},
		{"SHA:\t%v\n", s.Metadata.BuildCommit},
		{"Build Date:\t%v\n", s.Metadata.BuildDate},
	} {
		writeToConsole(v[0].(string), v[1])
	}
	writeToConsole("\n")

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

	go func() {
		fmt.Printf("\U0001F7E2 Server Started Listening on %s\n", s.server.Addr)
		serr := s.server.ListenAndServe()
		if serr != nil && serr != http.ErrServerClosed && err == nil {
			err = serr
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	fmt.Println("🛑 Server Stopped")

	err = s.gracefulShutdown()
	if err != nil {
		return err
	}

	if err := s.server.Close(); err != nil {
		logrus.Error("error while closing server:", err)
	}

	fmt.Println("💯 Server Exited Properly")
	return
}

func (s *Server) gracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}

// AddRoutes takes routes to register in server.
func (s *Server) AddRoutes(routes ...api.Route) {
	s.HTTPRoutes = append(s.HTTPRoutes, routes...)
}
