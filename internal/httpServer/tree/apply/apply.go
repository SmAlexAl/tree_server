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
}

type cacheStorage interface {
	GetAllTransaction() []model.Transaction
	Set(object model.Object)
	GetCollection() map[string]model.Object
	InvalidateTransaction()
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
				err := sqlStorage.SaveLeaf(transaction.Object)

				if err != nil {
					render.JSON(w, r, resp.Error(fmt.Errorf("add error: %s", err).Error()))
					return
				}
				break
			}
		}

		cache.InvalidateTransaction()
		tree, err := sqlStorage.GetTree()
		//TODO можно вынести чтобы убрать дубляж
		for _, object := range tree {
			if object.Active {
				object.State = model.ACTIVE_STATE
			} else {
				object.State = model.DELETE_STATE
			}

			tree[object.Id] = object
		}

		if err != nil {
			render.JSON(w, r, resp.Error(fmt.Errorf("select tree error: %s", err).Error()))
			return
		}

		render.JSON(w, r, resp.OK(viewer.GetData(tree)))
	}
}
