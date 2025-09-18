package response

const (
	StatusOK    = "ok"
	StatusError = "error"
)

type Response struct {
	Status string `json:"status"`
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func OK(result any) Response {
	return Response{
		Status: StatusOK,
		Result: result,
	}
}
