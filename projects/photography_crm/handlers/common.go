package handlers

type ActionStatusResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type GenericStatusResponse struct {
	Status string `json:"status"`
}

type ActionStatusSuccessResponse struct {
	Status string `json:"status"`
}

type GenericItemDeleteResponse struct {
	Status string `json:"status"`
	ID     string `json:"id"`
}

type DetailedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponseEnvelope struct {
	Error DetailedError `json:"error"`
}
