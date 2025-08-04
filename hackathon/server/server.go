package server

import (
	"elotus.com/hackathon/storage"
	"go.uber.org/zap"

	"elotus.com/hackathon/pkg/logutil"
)

type Server struct {
	Cfg     *Config
	Storage storage.Storage
	Logger  *zap.Logger
}

type Config struct {
	ServerPort                    int    `json:"server-port"`
	MysqlDatasource               string `json:"mysql-datasource"`
	MysqlConnMaxLifetimeInSeconds int    `json:"mysql-conn-max-life-time-in-seconds"`
	MysqlMaxOpenConns             int    `json:"mysql-max-open-conns"`
	MysqlMaxIdleConns             int    `json:"mysql-max-idle-conns"`
	AuthSecretKey                 string `json:"auth-secret-key"`
	TokenExpiresInSeconds         int    `json:"auth-token-expires-in-seconds"`
	UploadFileDir                 string `json:"upload-file-dir"`
}

func NewServer(cfg *Config) (*Server, error) {
	lg, err := logutil.CreateDefaultZapLogger(zap.InfoLevel)
	if err != nil {
		return nil, err
	}
	storageCfg := newStorageConfig(cfg)
	server := &Server{
		Storage: storage.NewStorage(lg, storageCfg),
		Cfg:     cfg,
		Logger:  lg,
	}

	return server, nil
}

func newStorageConfig(cfg *Config) *storage.Config {
	return &storage.Config{
		Datasource:               cfg.MysqlDatasource,
		ConnMaxLifeTimeInSeconds: cfg.MysqlConnMaxLifetimeInSeconds,
		MaxOpenConns:             cfg.MysqlMaxOpenConns,
		MaxIdleConns:             cfg.MysqlMaxIdleConns,
	}
}
