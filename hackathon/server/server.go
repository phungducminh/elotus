package server

import "elotus.com/hackathon/storage"

type Server struct {
	Cfg     *ServerConfig
	Storage storage.Storage
}

type ServerConfig struct {
	MysqlDatasource string `json:"mysql-datasource"`
	AuthSecretKey   string `json:"auth-secret-key"`
}

func NewServer(cfg *ServerConfig) *Server {
	storageCfg := &storage.Config{
		Datasource: cfg.MysqlDatasource,
	}

	return &Server{
		Storage: storage.NewStorage(storageCfg),
		Cfg:     cfg,
	}
}
