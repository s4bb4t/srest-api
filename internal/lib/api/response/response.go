package response

type Response struct {
	Error string `json:"msg,omitempty"`
}

func OK() Response {
	return Response{}
}

func Error(msg string) Response {
	return Response{
		Error: msg,
	}
}
