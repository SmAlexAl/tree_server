package response

type Response struct {
	Status     string      `json:"status"`
	Error      string      `json:"error,omitempty"`
	Collection interface{} `json:"collection,omitempty"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(c interface{}) Response {
	return Response{
		Status:     StatusOK,
		Collection: c,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
