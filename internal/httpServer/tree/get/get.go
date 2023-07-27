package get

import (
	"fmt"
	resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	"github.com/SmAlexAl/tree_server.git/internal/service"
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
	GetCollectionIndex() map[string][]string
}

func New(cache cacheStorage, sqlStorage sqlStorage, viewer viewer.Viewer) http.HandlerFunc {
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
			render.JSON(w, r, resp.Error(fmt.Errorf("id not found").Error()))

			return
		}

		_, ok := cache.Get(id)

		if !ok {
			object, err := sqlStorage.GetLeaf(id)

			if err != nil {
				render.JSON(w, r, resp.Error(err.Error()))

				return
			}

			updateCache(cache, object)
		}

		render.JSON(w, r, resp.OK(viewer.GetData(cache.GetCollection())))
	}
}

func updateCache(cache cacheStorage, newObject model.Object) {
	if newObject.Parent != "" {
		parentOb, ok := cache.Get(newObject.Parent)
		if ok {
			newObject.Active = parentOb.Active
		}
	}

	var res []string
	colIndex := cache.GetCollectionIndex()

	childrenId := service.GetChildren(res, newObject.Id, colIndex)

	currentActive := newObject.Active

	for _, id := range childrenId {
		child, _ := cache.Get(id)

		if currentActive && !child.Active {
			currentActive = child.Active
		}

		child.Active = currentActive

		cache.Set(child)
	}

	cache.Set(newObject)
}
