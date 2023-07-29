package apply

import (
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
	GetLeafsByActive(id string, active bool) (bool, error)
	BeginTransaction() error
	Commit() error
	Rollback() error
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
		err := sqlStorage.BeginTransaction()

		if err != nil {
			render.JSON(w, r, resp.Error(err.Error()))
		}

		for _, transaction := range transactionList {
			switch transaction.Command {
			case model.UPDATE:
				err := sqlStorage.UpdateLeaf(transaction.Object)

				if err != nil {
					render.JSON(w, r, resp.Error(err.Error()))

					sqlStorage.Rollback()
					return
				}

				break
			case model.DELETE:
				err := sqlStorage.DeleteLeaf(transaction.Object)

				if err != nil {
					render.JSON(w, r, resp.Error(err.Error()))

					sqlStorage.Rollback()
					return
				}

				break
			case model.ADD:
				status, err := sqlStorage.GetLeafsByActive(transaction.Object.Parent, false)
				if err != nil {
					render.JSON(w, r, resp.Error(err.Error()))

					sqlStorage.Rollback()
				} else if !status {
					err = sqlStorage.SaveLeaf(transaction.Object)

					if err != nil {
						render.JSON(w, r, resp.Error(err.Error()))

						sqlStorage.Rollback()
						return
					}
				}

				break
			}
		}

		err = sqlStorage.Commit()

		if err != nil {
			render.JSON(w, r, resp.Error(err.Error()))
		}

		cache.InvalidateTransaction()

		tree, err := sqlStorage.GetTree()

		if err != nil {
			render.JSON(w, r, resp.Error(err.Error()))
			return
		}

		cacheCollection := cache.GetCollection()
		newCache := make(map[string]model.Object)

		//перезапись текущего кэша, нужно для тестирования
		for _, object := range tree {
			if _, ok := cacheCollection[object.Id]; ok {
				newCache[object.Id] = object
			}
		}

		cache.SetCollection(newCache)

		render.JSON(w, r, OKWithDb(viewer.GetData(newCache), viewer.GetData(tree)))
	}
}
