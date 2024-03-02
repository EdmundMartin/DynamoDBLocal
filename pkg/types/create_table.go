package types

/*
map[AttributeDefinitions:[map[AttributeName:year AttributeType:N] map[AttributeName:title AttributeType:S]] KeySchema:[map[AttributeName:year KeyType:HASH] map[AttributeName:title KeyType:RANGE]] TableName:Hello]
*/

type AttributeDefinition struct {
	AttributeName string `json:"AttributeName"`
	AttributeType string `json:"AttributeType"`
}

type KeySchemaElement struct {
	AttributeName string `json:"AttributeName"`
	KeyType       string `json:"KeyType"`
}

type ProvisionedThroughput struct {
	ReadCapacityUnits  int64 `json:"ReadCapacityUnits"`
	WriteCapacityUnits int64 `json:"WriteCapacityUnits"`
}

// CreateTablePayload has the minimum types required in order to create a table
type CreateTablePayload struct {
	TableName             string                `json:"TableName"`
	AttributeDefinitions  []AttributeDefinition `json:"AttributeDefinitions"`
	KeySchema             []KeySchemaElement    `json:"KeySchema"`
	ProvisionedThroughput ProvisionedThroughput `json:"ProvisionedThroughput"`
}

// 400 Bad Request response
type ErrorBadRequest struct {
	Message string
}

type ProvisionedThroughPutResponse struct {
	LastDecreaseDateTime   int64 `json:"LastDecreaseDateTime"`
	LastIncreaseDateTime   int64 `json:"LastIncreaseDateTime"`
	NumberOfDecreasesToday int64 `json:"NumberOfDecreasesToday"`
	ReadCapacityUnits      int64 `json:"ReadCapacityUnits"`
	WriteCapacityUnits     int64 `json:"WriteCapacityUnits"`
}

type TableDescriptionResponse struct {
	AttributeDefinitions  []AttributeDefinition         `json:"AttributeDefinitions"`
	KeySchema             []KeySchemaElement            `json:"KeySchema"`
	CreationDateTime      int64                         `json:"CreationDateTime"`
	ItemCount             int64                         `json:"ItemCount"`
	ProvisionedThroughput ProvisionedThroughPutResponse `json:"ProvisionedThroughput"`
	TableArn              string                        `json:"TableArn"`
	TableName             string                        `json:"TableName"`
	TableStatus           string                        `json:"TableStatus"`
}

// CreateTableResponse has the response returned from Dynamo on table creation
type CreateTableResponse struct {
	TableDescription TableDescriptionResponse `json:"TableDescription"`
}
