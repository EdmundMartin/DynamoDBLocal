package main

import (
	"DynamoDBLocal/pkg/dummyclient"
	"DynamoDBLocal/pkg/types"
)

func main() {

	dummyclient.GetCreateTableResponse("http://localhost:8000", types.CreateTablePayload{
		TableName: "Hello",
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: "year",
				AttributeType: "N",
			},
			{
				AttributeName: "title",
				AttributeType: "S",
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				KeyType:       "HASH",
				AttributeName: "year",
			},
			{
				AttributeName: "title",
				KeyType:       "RANGE",
			},
		},
		ProvisionedThroughput: types.ProvisionedThroughput{
			ReadCapacityUnits:  10,
			WriteCapacityUnits: 10,
		},
	})
}
