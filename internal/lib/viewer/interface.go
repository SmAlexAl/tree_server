package viewer

import "github.com/SmAlexAl/tree_server.git/internal/model"

type Viewer interface {
	GetData(collection map[string]model.Object) interface{}
}
