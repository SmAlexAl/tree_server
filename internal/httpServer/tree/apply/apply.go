package apply

import (
	"fmt"
	resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	"github.com/go-chi/render"
	"net/http"
)

type sqlStorage interface {
	UpdateLeaf(object model.Object) error
	DeleteLeaf(object model.Object) error
	SaveLeaf(object model.Object) error
	GetTree() (map[string]model.Object, error)
	GetLeafsByActive(id string, active bool) error
}

type cacheStorage interface {
	GetAllTransaction() []model.Transaction
	Set(object model.Object)
	GetCollection() map[string]model.Object
	InvalidateTransaction()
	SetCollection(collection map[string]model.Object)
}

func New(cache cacheStorage, sqlStorage sqlStorage, viewer viewer.Viewer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactionList := cache.GetAllTransaction()

		for _, transaction := range transactionList {
			switch transaction.Command {
			case model.UPDATE:
				err := sqlStorage.UpdateLeaf(transaction.Object)

				if err != nil {
					render.JSON(w, r, resp.Error(fmt.Errorf("update error: %s", err).Error()))
					return
				}

				break
			case model.DELETE:
				err := sqlStorage.DeleteLeaf(transaction.Object)

				if err != nil {
					render.JSON(w, r, resp.Error(fmt.Errorf("delete error: %s", err).Error()))
					return
				}

				break
			case model.ADD:
				err := sqlStorage.GetLeafsByActive(transaction.Object.Parent, false)
				//TODO решение не оптимальное, по возможности переделать
				if err != nil && err.Error() == "sql: no rows in result set" {
					err = sqlStorage.SaveLeaf(transaction.Object)

					if err != nil {
						render.JSON(w, r, resp.Error(fmt.Errorf("add error: %s", err).Error()))
						return
					}
				} else if err != nil {
					render.JSON(w, r, resp.Error(fmt.Errorf("add error(select): %s", err).Error()))
				}

				break
			}
		}

		cache.InvalidateTransaction()
		tree, err := sqlStorage.GetTree()
		cacheCollection := cache.GetCollection()
		newCache := make(map[string]model.Object)
		//TODO можно вынести чтобы убрать дубляж
		for _, object := range tree {
			if object.Active {
				object.State = model.ACTIVE_STATE
			} else {
				object.State = model.DELETE_STATE
			}

			if _, ok := cacheCollection[object.Id]; ok {
				newCache[object.Id] = object
			}

			tree[object.Id] = object
		}

		cache.SetCollection(newCache)

		if err != nil {
			render.JSON(w, r, resp.Error(fmt.Errorf("select tree error: %s", err).Error()))
			return
		}

		render.JSON(w, r, OKWithDb(viewer.GetData(newCache), viewer.GetData(tree)))
	}
}
