package config

import (
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"sync"
)

var (
	once sync.Once
	conf *Config
)

// New initializes the configuration setting. It searches for the
// config file with a filename of "config.yaml".
//
// Config file will be search in the following path in order:
// - "/etc/noteapp"
// - "/run/secrets"
// - "."
func New() *Config {
	// Do singleton
	once.Do(func() {
		var err error
		conf, err = newConfig(afero.NewOsFs())
		if err != nil {
			panic(err)
		}
	})

	return conf
}

func newConfig(fs afero.Fs) (*Config, error) {
	viper.SetFs(fs)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("/etc/noteapp")
	// For Docker Compose default path for mounting the secret.
	// see. https://docs.docker.com/compose/compose-file/compose-file-v3/#secrets
	viper.AddConfigPath("/run/secrets")
	viper.AddConfigPath("/")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	if viper.Get("store.file.path") == nil {
		viper.Set("store.file.path", ".")
	}

	if viper.Get("server.port") == nil {
		viper.Set("server.port", 50001)
	}

	var conf Config
	err = viper.Unmarshal(&conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}

// Config is the application-level configuration containing
// all the information for running the application.
type Config struct {
	Server Server
	// Store Database Configuration
	Store Store
}

// Server contains the server configuration.
type Server struct {
	// Port is the port of the server when its value is empty
	// in config file the default "50001" will be use.
	Port int
}

// Store contains the store database configuration.
type Store struct {
	File File
}

// File contains the file store configuration.
type File struct {
	// path is the path where the files of the file store will be store.
	// When its value is empty in config file the default "." will be use.
	Path string
}
