package server

import (
	"fmt"
	"log"
	"net/http"

	"elotus.com/hackathon/storage"
	"go.uber.org/zap"

	"elotus.com/hackathon/pkg/logutil"
)

var defaultOptions = &Options{
	AuthSecretKey:                 "SECRET-KEY",
	TokenExpiresInSeconds:         60,
	MysqlDatasource:               "root:elotus@tcp(localhost:3306)/elotus",
	MysqlConnMaxLifetimeInSeconds: 0, // conn are not closed due to a connection's age
	MysqlMaxOpenConns:             1000,
	MysqlMaxIdleConns:             200,
	ServerPort:                    8080,
	UploadFileDir:                 "tmp",
}

type Server struct {
	Cfg     *Options
	Storage storage.Storage
	Logger  *zap.Logger
}

type Options struct {
	ServerPort                    int    `json:"server-port"`
	MysqlDatasource               string `json:"mysql-datasource"`
	MysqlConnMaxLifetimeInSeconds int    `json:"mysql-conn-max-life-time-in-seconds"`
	MysqlMaxOpenConns             int    `json:"mysql-max-open-conns"`
	MysqlMaxIdleConns             int    `json:"mysql-max-idle-conns"`
	AuthSecretKey                 string `json:"auth-secret-key"`
	TokenExpiresInSeconds         int    `json:"auth-token-expires-in-seconds"`
	UploadFileDir                 string `json:"upload-file-dir"`
}

type Option func(*Options)

func WithUploadFileDir(dir string) Option {
	return Option(func(opts *Options) {
		opts.UploadFileDir = dir
	})
}

func NewServer(opt ...Option) (*Server, error) {
	opts := defaultOptions
	for _, o := range opt {
		o(opts)
	}
	lg, err := logutil.CreateDefaultZapLogger(zap.InfoLevel)
	if err != nil {
		return nil, err
	}
	storageCfg := newStorageConfig(opts)
	server := &Server{
		Storage: storage.NewStorage(lg, storageCfg),
		Cfg:     opts,
		Logger:  lg,
	}

	return server, nil
}

func (s *Server) Serve() error {
	log.Printf("Starting server on %d...", s.Cfg.ServerPort)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Cfg.ServerPort), nil)
}

func newStorageConfig(opts *Options) *storage.Config {
	return &storage.Config{
		Datasource:               opts.MysqlDatasource,
		ConnMaxLifeTimeInSeconds: opts.MysqlConnMaxLifetimeInSeconds,
		MaxOpenConns:             opts.MysqlMaxOpenConns,
		MaxIdleConns:             opts.MysqlMaxIdleConns,
	}
}
