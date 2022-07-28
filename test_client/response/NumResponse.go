package response

type NumResponse struct {
	Result int `json:"result"`
	ErrorMsg string `json:"error"`
}
