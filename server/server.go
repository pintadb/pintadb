package server

import (
	"os"
	"path/filepath"

	"github.com/columbusearch/pintadb/db"
)

type PintaDBServer struct {
	Documents []db.TextVec
	Config    Config
}

type Config struct {
	FullPath  string
	Dimension int
}

// NewServer creates a new PintaDB server
// The server is initialized with an empty database by default
func NewServer(config Config) (*PintaDBServer, error) {
	server := &PintaDBServer{
		Config: config,
	}

	err := os.MkdirAll(filepath.Dir(config.FullPath), os.ModePerm)
	if err != nil {
		return nil, err
	}

	file, err := os.Create(server.Config.FullPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return server, nil
}
