package appTypes

type WSResp struct {
	StatusCode int    `json:"statusCode"`
	Body       any    `json:"body"`
	Error      string `json:"error"`
}
