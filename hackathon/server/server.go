package server

import "elotus.com/hackathon/storage"

type Server struct {
	Cfg     *Config
	Storage storage.Storage
}

type Config struct {
	MysqlDatasource               string `json:"mysql-datasource"`
	MysqlConnMaxLifetimeInSeconds int    `json:"mysql-conn-max-life-time-in-seconds"`
	MysqlMaxOpenConns             int    `json:"mysql-max-open-conns"`
	MysqlMaxIdleConns             int    `json:"mysql-max-idle-conns"`
	AuthSecretKey                 string `json:"auth-secret-key"`
}

func NewServer(cfg *Config) *Server {
	storageCfg := newStorageConfig(cfg)
	return &Server{
		Storage: storage.NewStorage(storageCfg),
		Cfg:     cfg,
	}
}

func newStorageConfig(cfg *Config) *storage.Config {
	return &storage.Config{
		Datasource:               cfg.MysqlDatasource,
		ConnMaxLifeTimeInSeconds: cfg.MysqlConnMaxLifetimeInSeconds,
		MaxOpenConns:             cfg.MysqlMaxOpenConns,
		MaxIdleConns:             cfg.MysqlMaxIdleConns,
	}
}
