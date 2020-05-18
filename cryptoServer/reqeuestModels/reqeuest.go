package requestModels

type Request struct {
	RequestType int    `json:"requestType"`
	Data        []byte `json:"data"`
}
