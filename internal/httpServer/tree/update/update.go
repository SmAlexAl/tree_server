package update

import (
	"fmt"
	resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	"github.com/go-chi/render"
	"net/http"
)

type Request struct {
	Id    string
	Value string
}

type cacheStorage interface {
	Get(id string) (model.Object, bool)
	Set(value model.Object)
	GetCollection() map[string]model.Object
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

		val.Value = req.Value
		cache.Set(val)

		cache.AddTransaction(model.NewTransaction(model.UPDATE, val))

		render.JSON(w, r, resp.OK(viewer.GetData(cache.GetCollection())))
	}
}
