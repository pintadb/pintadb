package server

import (
	"encoding/gob"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sync"

	"github.com/columbusearch/pintadb/db"
)

type PintaDBServer struct {
	Documents []db.TextVec
	Config    Config
}

type Config struct {
	FullPath  string
	Dimension uint64
	HTTPPort  uint64
	GRPCPort  uint64
}

// NewServer creates a new PintaDB server
// The server is initialized with an empty database by default
func NewServer(config Config) (*PintaDBServer, error) {
	server := &PintaDBServer{
		Config: config,
	}

	if _, err := os.Stat(config.FullPath); os.IsNotExist(err) {
		err = createDBFiles(config.FullPath)
		if err != nil {
			return nil, err
		}
	} else {
		server.LoadDB()
	}

	fmt.Printf("HTTP Server is running on port %d", config.HTTPPort)
	NewHTTPServer(server, config.HTTPPort)
	// NewGRPCServer(server, config.GRPCPort)

	return server, nil
}

func createDBFiles(path string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

// LoadDB loads the database from the file
func (s *PintaDBServer) LoadDB() error {
	file, err := s.DBFile()
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	fmt.Println("Loading data from ", s.Config.FullPath)
	// decoder.Decode(&s.Documents)

	err = decoder.Decode(&s.Documents)
	if err != nil {
		return err
	}

	return nil
}

// DBFile returns the file object of the database
func (s *PintaDBServer) DBFile() (*os.File, error) {
	file, err := os.OpenFile(s.Config.FullPath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// AddDocument adds a new document to the database
// The document is indexed and vectorized before being added
func (s *PintaDBServer) AddDocument(text string) error {
	t := db.TextVec{WordFreq: db.MakeWordFreq(text), Vector: make([]float64, s.Config.Dimension), RawText: text}
	t.CalculateVector(s.Config.Dimension)

	s.LoadDB()
	flushFile(s.Config.FullPath)

	s.Documents = append(s.Documents, t)

	file, err := s.DBFile()
	if err != nil {
		return err
	}
	defer file.Close()

	err = encodeToFile(file, s.Documents)
	if err != nil {
		return err
	}

	s.Documents = []db.TextVec{}
	return nil
}

func flushFile(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func encodeToFile(file *os.File, vectors []db.TextVec) error {
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(vectors)
	return err
}

// Query returns the k most similar documents to the given text
// And the returned documents are sorted by similarity in descending order
func (s *PintaDBServer) Search(text string, k int) ([]int, error) {
	queryVec := db.TextVec{WordFreq: db.MakeWordFreq(text), Vector: make([]float64, s.Config.Dimension)}
	queryVec.CalculateVector(s.Config.Dimension)
	s.LoadDB()
	return calculateResult(queryVec, s.Documents, k), nil
}

func calculateResult(queryVec db.TextVec, documents []db.TextVec, k int) []int {
	result := make([]int, k)
	resultDist := make([]float64, k)
	for i := range result {
		result[i] = -1
		resultDist[i] = math.Inf(1)
	}

	type resultPair struct {
		index int
		dist  float64
	}
	results := make(chan resultPair, len(documents))

	var wg sync.WaitGroup
	for i, doc := range documents {
		wg.Add(1)
		go func(i int, doc db.TextVec) {
			defer wg.Done()
			dist := db.CosineDistance(doc.Vector, queryVec.Vector)
			results <- resultPair{i, dist}
		}(i, doc)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for pair := range results {
		if pair.dist >= resultDist[k-1] {
			continue
		}
		j := k - 2
		for ; j >= 0; j-- {
			if pair.dist < resultDist[j] {
				result[j+1], resultDist[j+1] = result[j], resultDist[j]
			} else {
				break
			}
		}
		result[j+1], resultDist[j+1] = pair.index, pair.dist
	}

	return result
}
