package reset

import (
	resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"
	"github.com/SmAlexAl/tree_server.git/internal/lib/viewer"
	"github.com/SmAlexAl/tree_server.git/internal/model"
	"github.com/go-chi/render"
	"net/http"
)

type sqlStorage interface {
	TruncateTree() error
	SaveLeafs(fixtures []model.Object) error
	GetTree() (map[string]model.Object, error)
}

type cacheStorage interface {
	Invalidate()
	GetCollection() map[string]model.Object
}

type fixtures interface {
	GetCollection() []model.Object
}

func New(cache cacheStorage, sqlStorage sqlStorage, fixtures fixtures, viewer viewer.Viewer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := sqlStorage.TruncateTree()

		if err != nil {
			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		err = sqlStorage.SaveLeafs(fixtures.GetCollection())

		if err != nil {
			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		cache.Invalidate()

		collection, err := sqlStorage.GetTree()

		if err != nil {
			render.JSON(w, r, resp.Error(err.Error()))

			return
		}

		render.JSON(w, r, resp.OK(viewer.GetData(collection)))
	}
}
