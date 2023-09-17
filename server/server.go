package server

import (
	"encoding/gob"
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

func (s *PintaDBServer) DBFile() (*os.File, error) {
	file, err := os.OpenFile(s.Config.FullPath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *PintaDBServer) AddDocument(text string) error {
	t := db.TextVec{WordFreq: db.MakeWordFreq(text), Vector: make([]float64, s.Config.Dimension), RawText: text}
	t.CalculateVector(s.Config.Dimension)
	s.Documents = append(s.Documents, t)

	file, err := s.DBFile()
	if err != nil {
		return err
	}
	defer file.Close()

	err = encode(file, s.Documents)
	if err != nil {
		return err
	}

	s.Documents = []db.TextVec{}
	return nil
}

func encode(file *os.File, vectors []db.TextVec) error {
	encoder := gob.NewEncoder(file)
	err := encoder.Encode(vectors)
	return err
}

// Query returns the k most similar documents to the given text
// And the returned documents are sorted by similarity in descending order
func (s *PintaDBServer) Query(text string, k int) ([]int, error) {
	queryVec := db.TextVec{WordFreq: db.MakeWordFreq(text), Vector: make([]float64, s.Config.Dimension)}
	queryVec.CalculateVector(s.Config.Dimension)

	file, err := s.DBFile()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&s.Documents)
	if err != nil {
		return nil, err
	}

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
