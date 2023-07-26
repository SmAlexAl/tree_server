package apply

import resp "github.com/SmAlexAl/tree_server.git/internal/lib/api/response"

type Response struct {
	resp.Response
	DB interface{} `json:"db,omitempty"`
}

func OKWithDb(collection, db interface{}) Response {
	return Response{
		Response: resp.OK(collection),
		DB:       db,
	}
}
