package get

import (
	"fmt"
	resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	"github.com/go-chi/render"
	"net/http"
)

type sqlStorage interface {
	GetLeaf(id string) (model.Object, error)
}

type cacheStorage interface {
	Get(id string) (model.Object, bool)
	Set(object model.Object)
	GetCollection() map[string]model.Object
}

func New(storage cacheStorage, sqlStorage sqlStorage, viewer viewer.Viewer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		var id string
		if _, ok := q["id"]; ok {
			id = q["id"][0]
		} else {
			render.JSON(w, r, resp.Error(fmt.Errorf("request parse error").Error()))

			return
		}

		if id == "" {
			render.JSON(w, r, resp.Error(fmt.Errorf("request parse error").Error()))

			return
		}

		_, ok := storage.Get(id)

		if !ok {
			object, err := sqlStorage.GetLeaf(id)

			if err != nil {
				render.JSON(w, r, resp.Error(fmt.Errorf("select data error: %s", err).Error()))

				return
			}

			storage.Set(object)
		}

		render.JSON(w, r, resp.OK(viewer.GetData(storage.GetCollection())))
	}
}
