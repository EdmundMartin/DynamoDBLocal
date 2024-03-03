package types

type Capacity struct {

	// The total number of capacity units consumed on a table or an index.
	CapacityUnits float64

	// The total number of read capacity units consumed on a table or an index.
	ReadCapacityUnits float64

	// The total number of write capacity units consumed on a table or an index.
	WriteCapacityUnits float64
}

type ConsumedCapacity struct {
	CapacityUnits float64 `json:"CapacityUnits"`

	// The amount of throughput consumed on each global index affected by the
	// operation.
	GlobalSecondaryIndexes map[string]Capacity `json:"GlobalSecondaryIndexes"`

	// The amount of throughput consumed on each local index affected by the operation.
	LocalSecondaryIndexes map[string]Capacity `json:"LocalSecondaryIndexes"`

	// The total number of read capacity units consumed by the operation.
	ReadCapacityUnits float64 `json:"ReadCapacityUnits"`

	// The amount of throughput consumed on the table affected by the operation.
	Table Capacity `json:"Table"`

	// The name of the table that was affected by the operation.
	TableName string `json:"TableName"`

	// The total number of write capacity units consumed by the operation.
	WriteCapacityUnits float64 `json:"WriteCapacityUnits"`
}

type PutItemRequest struct {
	Item      map[string]map[string]string `json:"Item"`
	TableName string                       `json:"TableName"`
}

type PutItemResponse struct {
	// TODO - Attributes requires support of return values in request
	// TODO - ItemCollectionMetrics requires support of ReturnItemCollectionMetrics in request
	ConsumedCapacity ConsumedCapacity `json:"ConsumedCapacity"`
}

// NewDefaultWriteConsumedCapacity returns a ConsumedCapacity reflecting the fact that we don't actually use any
// write capacity
func NewDefaultWriteConsumedCapacity(tableName string) ConsumedCapacity {
	return ConsumedCapacity{
		CapacityUnits:     1.0,
		ReadCapacityUnits: 0,
		Table: Capacity{
			CapacityUnits:      1.0,
			WriteCapacityUnits: 1,
		},
		TableName:          tableName,
		WriteCapacityUnits: 1,
	}
}
