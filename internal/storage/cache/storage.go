package cache

import (
	"github.com/SmAlexAl/tree_server.git/internal/model"
)

type Storage struct {
	collection  map[string]model.Object
	index       map[string][]string
	transaction []model.Transaction
}

func New() (*Storage, error) {
	collection := make(map[string]model.Object)
	index := make(map[string][]string)
	return &Storage{
		collection: collection,
		index:      index,
	}, nil
}

func (s *Storage) Get(id string) (model.Object, bool) {
	obj, ok := s.collection[id]

	return obj, ok
}

func (s *Storage) Set(value model.Object) {
	if _, ok := s.collection[value.Id]; !ok {
		s.index[value.Parent] = append(s.index[value.Parent], value.Id)
	}
	s.collection[value.Id] = value
}

func (s *Storage) SetCollection(collection map[string]model.Object) {
	s.collection = collection

	s.index = map[string][]string{}
	for _, value := range s.collection {
		if value.Active {
			s.index[value.Parent] = append(s.index[value.Parent], value.Id)
		}
	}
}

func (s *Storage) Invalidate() {
	s.collection = make(map[string]model.Object)
	s.index = make(map[string][]string)
	s.transaction = []model.Transaction{}
}

func (s *Storage) InvalidateTransaction() {
	s.transaction = []model.Transaction{}
}

func (s *Storage) GetCollection() map[string]model.Object {
	return s.collection
}

func (s *Storage) GetCollectionIndex() map[string][]string {
	return s.index
}

func (s *Storage) AddTransaction(val model.Transaction) {
	s.transaction = append(s.transaction, val)
}

func (s *Storage) GetAllTransaction() []model.Transaction {
	return s.transaction
}
