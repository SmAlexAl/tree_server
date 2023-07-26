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
	for _, val := range collection {
		if val.Active == false {
			continue
		}
		parent := val.Parent
		if _, ok := collection[parent]; !ok || parent == "" {
			res = append(res, Object{Id: val.Id, Parent: "#", Text: val.Value})

		} else {
			res = append(res, Object{Id: val.Id, Parent: parent, Text: val.Value})
		}
	}
	return res
}
