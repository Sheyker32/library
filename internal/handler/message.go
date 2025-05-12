package handler

type Response struct {
	Success   bool        `json:"success"`
	ErrorCode int         `json:"error_code,omitempty"`
	Data      interface{} `json:"data"`
}

type Data struct {
	Message string `json:"message"`
}
