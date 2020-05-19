package requestModels

type Request struct {
	RequestType int    `json:"requestType"`
	UserID      string `json:"userID"`
	Data        []byte `json:"data"`
}
