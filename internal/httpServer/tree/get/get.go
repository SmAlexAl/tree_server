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
			render.JSON(w, r, resp.Error(fmt.Errorf("request parse error").Error()))

			return
		}

		_, ok := cache.Get(id)

		if !ok {
			object, err := sqlStorage.GetLeaf(id)

			if err != nil {
				render.JSON(w, r, resp.Error(fmt.Errorf("select data error: %s", err).Error()))

				return
			}

			updateCache(cache, object)
		}

		render.JSON(w, r, resp.OK(viewer.GetData(cache.GetCollection())))
	}
}

// TODO можно вынести в отдельный сервис
func updateCache(cache cacheStorage, newObject model.Object) {
	if newObject.Parent == "" {
		if newObject.Active {
			newObject.State = model.ACTIVE_STATE
		} else {
			newObject.State = model.DELETE_STATE
		}
	} else {
		parentOb, ok := cache.Get(newObject.Parent)

		if !newObject.Active {
			newObject.State = model.DELETE_STATE
		} else if ok {
			newObject.Active = parentOb.Active
			newObject.State = parentOb.State
		}
	}

	var res []string
	colIndex := cache.GetCollectionIndex()

	childrenId := getChildren(res, newObject.Id, colIndex)

	currentActive := newObject.Active
	currentState := newObject.State

	for _, v := range childrenId {
		child, _ := cache.Get(v)

		if currentActive && !child.Active {
			currentActive = child.Active
			currentState = child.State
		}

		child.Active = currentActive
		child.State = currentState

		cache.Set(child)
	}

	cache.Set(newObject)
}

func getChildren(res []string, id string, tree map[string][]string) []string {
	for _, val := range tree[id] {
		if _, ok := tree[val]; ok {
			res = getChildren(res, val, tree)
		}
		res = append(res, val)
	}

	return res
}
