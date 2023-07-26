package add

import (
	"fmt"
	resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	"github.com/go-chi/render"
	"net/http"
)

type Request struct {
	Value    string
	ParentId string
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

		parentObject, ok := cache.Get(req.ParentId)
		if !ok || !parentObject.Active || parentObject.State == model.UNKNOW_STATE {
			render.JSON(w, r, resp.Error(fmt.Errorf("parent object not found").Error()))

			return
		}

		val := model.NewObject(req.Value, req.ParentId)

		cache.Set(val)
		cache.AddTransaction(model.NewTransaction(model.ADD, val))

		render.JSON(w, r, resp.OK(viewer.GetData(cache.GetCollection())))
	}
}
