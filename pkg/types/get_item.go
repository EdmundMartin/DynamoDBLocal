package types

// TODO - Absolute minimal request
type GetItemRequest struct {
	TableName string                       `json:"TableName"`
	Key       map[string]map[string]string `json:"Item"`
}

type GetItemResponse struct {
	Item             map[string]map[string]string `json:"Item"` // TODO - Should potentially use attribute values here
	ConsumedCapacity ConsumedCapacity             `json:"ConsumedCapacity"`
}
