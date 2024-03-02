package types

type PutItemRequest struct {
	Item      map[string]map[string]string `json:"Item"`
	TableName string                       `json:"TableName"`
}
