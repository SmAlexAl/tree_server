package delete

import (
	"fmt"
	resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	"github.com/SmAlexAl/tree_server.git/internal/service"
	"github.com/go-chi/render"
	"net/http"
)

type Request struct {
	Id string
}

type cacheStorage interface {
	Get(id string) (model.Object, bool)
	GetCollectionIndex() map[string][]string
	GetCollection() map[string]model.Object
	SetCollection(collection map[string]model.Object)
	AddTransaction(transaction model.Transaction)
}

func New(cache cacheStorage, viewer viewer.Viewer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			render.JSON(w, r, resp.Error(fmt.Errorf("request parse error").Error()))

			return
		}

		val, ok := cache.Get(req.Id)

		if !ok {
			render.JSON(w, r, resp.Error(fmt.Errorf("object not found: %s", req.Id).Error()))

			return
		}

		index := cache.GetCollectionIndex()

		deleteIds := []string{val.Id}

		deleteIds = service.GetChildren(deleteIds, val.Id, index)

		cache.SetCollection(updateCache(cache.GetCollection(), deleteIds))

		cache.AddTransaction(model.NewTransaction(model.DELETE, val))

		render.JSON(w, r, resp.OK(viewer.GetData(cache.GetCollection())))
	}
}

func updateCache(collection map[string]model.Object, ids []string) map[string]model.Object {
	for _, id := range ids {
		v := collection[id]
		v.Active = false
		collection[id] = v
	}

	return collection
}
