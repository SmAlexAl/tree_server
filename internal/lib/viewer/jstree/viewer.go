package jstree

import (
	"github.com/SmAlexAl/tree_server.git/internal/model"
)

type Viewer struct {
}

type Object struct {
	Id     string `json:"id"`
	Parent string `json:"parent"`
	Text   string `json:"text"`
}

func New() *Viewer {
	return &Viewer{}
}

func (v Viewer) GetData(collection map[string]model.Object) interface{} {
	res := make([]Object, 0, len(collection))
	for _, object := range collection {
		value := object.Value + object.State
		parent := object.Parent
		_, ok := collection[parent]

		if !ok || parent == "" {
			res = append(res, Object{Id: object.Id, Parent: "#", Text: value})
		} else {
			res = append(res, Object{Id: object.Id, Parent: parent, Text: value})
		}
	}
	return res
}
