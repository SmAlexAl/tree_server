package initDb

import (
	"fmt"
	resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	"github.com/go-chi/render"
	"net/http"
)

type sqlStorage interface {
	GetTree() (map[string]model.Object, error)
}

type cacheStorage interface {
	Set(object model.Object)
	GetCollection() map[string]model.Object
}

func New(cache cacheStorage, sqlStorage sqlStorage, viewer viewer.Viewer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := sqlStorage.GetTree()
		if err != nil {
			render.JSON(w, r, resp.Error(fmt.Errorf("db error").Error()))
		}

		render.JSON(w, r, OKWithDb(viewer.GetData(cache.GetCollection()), viewer.GetData(res)))
	}
}
